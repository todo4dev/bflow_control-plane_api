package pgx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"src/application/adapter/database"
	"src/core/builder"
	"src/core/common"
)

// #region queryRunner

type queryRunner interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// #endregion
// #region pgxExecutor (helper interno, NÃO implementa interface de adapter)

type pgxExecutor struct {
	config *database.DatabaseConfig
	runner queryRunner
}

func quoteTable(table string) string {
	parts := strings.Split(table, ".")
	identifier := pgx.Identifier(parts)
	return identifier.Sanitize()
}

func escapeJSONField(field string) string {
	return strings.ReplaceAll(field, "'", "''")
}

func castJSONField(baseExpr string, goType reflect.Type) string {
	if goType == nil {
		return baseExpr
	}

	switch goType.Kind() {
	case reflect.Bool:
		return fmt.Sprintf("(%s)::boolean", baseExpr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64:
		return fmt.Sprintf("(%s)::numeric", baseExpr)
	case reflect.Struct:
		if goType.AssignableTo(reflect.TypeOf(time.Time{})) {
			return fmt.Sprintf("(%s)::timestamptz", baseExpr)
		}
	}

	// default: texto
	return baseExpr
}

func sqlOperator(op builder.WhereEnum) (string, error) {
	switch op {
	case builder.WhereEnum_Equal:
		return "=", nil
	case builder.WhereEnum_NotEqual:
		return "<>", nil
	case builder.WhereEnum_Like:
		return "LIKE", nil
	case builder.WhereEnum_NotLike:
		return "NOT LIKE", nil
	case builder.WhereEnum_GreaterThan:
		return ">", nil
	case builder.WhereEnum_NotGreaterThan:
		return "<=", nil
	case builder.WhereEnum_GreaterEqual:
		return ">=", nil
	case builder.WhereEnum_NotGreaterEqual:
		return "<", nil
	case builder.WhereEnum_LowerThan:
		return "<", nil
	case builder.WhereEnum_NotLowerThan:
		return ">=", nil
	case builder.WhereEnum_LowerEqual:
		return "<=", nil
	case builder.WhereEnum_NotLowerEqual:
		return ">", nil
	default:
		return "", fmt.Errorf("unsupported where operator: %s", op)
	}
}

func (e *pgxExecutor) buildWhereAndArgs(text *string, where *builder.WherePointerMap, index int) (string, []any, error) {
	var parts []string
	var args []any
	idx := index

	addArg := func(v any) int {
		args = append(args, v)
		pos := idx
		idx++
		return pos
	}

	if text != nil && *text != "" {
		pos := addArg(*text)
		parts = append(parts, fmt.Sprintf("data::text ILIKE '%%' || $%d || '%%'", pos))
	}

	if where != nil {
		// ordenar campos só pra deixar determinístico
		fieldNames := make([]string, 0, len(*where))
		for field := range *where {
			fieldNames = append(fieldNames, field)
		}
		sort.Strings(fieldNames)

		for _, field := range fieldNames {
			ops := (*where)[field]

			opKeys := make([]builder.WhereEnum, 0, len(ops))
			for op := range ops {
				opKeys = append(opKeys, op)
			}
			sort.Slice(opKeys, func(i, j int) bool { return string(opKeys[i]) < string(opKeys[j]) })

			baseExpr := fmt.Sprintf("data->>'%s'", escapeJSONField(field))

			for _, op := range opKeys {
				val := ops[op]

				switch op {
				case builder.WhereEnum_Empty:
					parts = append(parts, fmt.Sprintf("COALESCE(%s, '') = ''", baseExpr))
				case builder.WhereEnum_NotEmpty:
					parts = append(parts, fmt.Sprintf("COALESCE(%s, '') <> ''", baseExpr))
				case builder.WhereEnum_In, builder.WhereEnum_NotIn:
					v := reflect.ValueOf(val)
					if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
						return "", nil, fmt.Errorf("IN/NIN expects slice/array for field %s", field)
					}
					elemType := v.Type().Elem()
					typedExpr := castJSONField(baseExpr, elemType)

					pos := addArg(val)

					if op == builder.WhereEnum_In {
						parts = append(parts, fmt.Sprintf("%s = ANY($%d)", typedExpr, pos))
					} else {
						parts = append(parts, fmt.Sprintf("NOT (%s = ANY($%d))", typedExpr, pos))
					}
				default:
					goType := reflect.TypeOf(val)
					typedExpr := castJSONField(baseExpr, goType)

					sqlOp, err := sqlOperator(op)
					if err != nil {
						return "", nil, err
					}

					pos := addArg(val)
					parts = append(parts, fmt.Sprintf("%s %s $%d", typedExpr, sqlOp, pos))
				}
			}
		}
	}

	if len(parts) == 0 {
		return "", args, nil
	}

	return strings.Join(parts, " AND "), args, nil
}

func (e *pgxExecutor) buildWhereFromQuery(query *builder.Query[json.RawMessage], index int) (string, []any, error) {
	if query == nil {
		return "", nil, nil
	}
	return e.buildWhereAndArgs(query.TextCond, query.WhereCond, index)
}

func (e *pgxExecutor) buildWhereFromWhere(where *builder.WhereBuilder[json.RawMessage], index int) (string, []any, error) {
	if where == nil {
		return "", nil, nil
	}
	pointerMap := where.PointerMap
	return e.buildWhereAndArgs(nil, &pointerMap, index)
}

func (e *pgxExecutor) buildOrderByClause(query *builder.Query[json.RawMessage]) string {
	if query == nil || query.SortCond == nil {
		return ""
	}

	sortMap := *query.SortCond
	if len(sortMap) == 0 {
		return ""
	}

	fieldNames := make([]string, 0, len(sortMap))
	for field := range sortMap {
		fieldNames = append(fieldNames, field)
	}
	sort.Strings(fieldNames)

	parts := make([]string, 0, len(fieldNames))
	for _, field := range fieldNames {
		dir := "ASC"
		if sortMap[field] == builder.SortEnum_Desc {
			dir = "DESC"
		}
		parts = append(parts, fmt.Sprintf("data->>'%s' %s", escapeJSONField(field), dir))
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, ", ")
}

func (e *pgxExecutor) FindOne(ctx context.Context, table string, query *builder.Query[json.RawMessage]) (*json.RawMessage, error) {
	where, args, err := e.buildWhereFromQuery(query, 1)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("SELECT data FROM %s", quoteTable(table))
	if where != "" {
		sql += " WHERE " + where
	}

	orderBy := e.buildOrderByClause(query)
	if orderBy != "" {
		sql += " ORDER BY " + orderBy
	}

	sql += " LIMIT 1"

	row := e.runner.QueryRow(ctx, sql, args...)

	var rawBytes []byte
	if err := row.Scan(&rawBytes); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	msg := json.RawMessage(rawBytes)
	return &msg, nil
}

func (e *pgxExecutor) FindMany(ctx context.Context, table string, query *builder.Query[json.RawMessage]) (*builder.Result[json.RawMessage], error) {
	if query == nil {
		query = builder.NewQuery[json.RawMessage]()
	}

	limit := e.config.DefaultLimit
	if query.LimitCond != nil && *query.LimitCond > 0 {
		limit = *query.LimitCond
	}
	if limit <= 0 {
		limit = 50
	}

	var offset int64
	if query.OffsetCond != nil && *query.OffsetCond >= 0 {
		offset = *query.OffsetCond
	}

	whereClause, args, err := e.buildWhereFromQuery(query, 1)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("SELECT data FROM %s", quoteTable(table))
	if whereClause != "" {
		sql += " WHERE " + whereClause
	}

	orderBy := e.buildOrderByClause(query)
	if orderBy != "" {
		sql += " ORDER BY " + orderBy
	}

	sql += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := e.runner.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]json.RawMessage, 0)
	for rows.Next() {
		var rawBytes []byte
		if err := rows.Scan(&rawBytes); err != nil {
			return nil, err
		}
		items = append(items, json.RawMessage(rawBytes))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	total, err := e.Count(ctx, table, query)
	if err != nil {
		return nil, err
	}

	return &builder.Result[json.RawMessage]{
		Offset: offset,
		Limit:  limit,
		Total:  total,
		Items:  items,
	}, nil
}

func (e *pgxExecutor) Count(ctx context.Context, table string, query *builder.Query[json.RawMessage]) (int64, error) {
	where, args, err := e.buildWhereFromQuery(query, 1)
	if err != nil {
		return 0, err
	}

	sql := fmt.Sprintf("SELECT COUNT(*) FROM %s", quoteTable(table))
	if where != "" {
		sql += " WHERE " + where
	}

	row := e.runner.QueryRow(ctx, sql, args...)

	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (e *pgxExecutor) Insert(ctx context.Context, table string, entities []json.RawMessage) error {
	if len(entities) == 0 {
		return nil
	}

	sql := fmt.Sprintf("INSERT INTO %s (data) VALUES ", quoteTable(table))

	args := make([]any, 0, len(entities))
	valuePlaceholders := make([]string, 0, len(entities))

	for i, entity := range entities {
		args = append(args, []byte(entity))
		valuePlaceholders = append(valuePlaceholders, fmt.Sprintf("($%d::jsonb)", i+1))
	}

	sql += strings.Join(valuePlaceholders, ", ")

	_, err := e.runner.Exec(ctx, sql, args...)
	return err
}

func (e *pgxExecutor) Update(ctx context.Context, table string, where *builder.WhereBuilder[json.RawMessage], update *builder.UpdateBuilder[json.RawMessage]) (int64, error) {
	if update == nil || len(update.Changes) == 0 {
		return 0, nil
	}

	patch, err := json.Marshal(update.Changes)
	if err != nil {
		return 0, err
	}

	// $1 é o patch; WHERE começa do $2
	whereClause, args, err := e.buildWhereFromWhere(where, 2)
	if err != nil {
		return 0, err
	}

	finalArgs := make([]any, 0, len(args)+1)
	finalArgs = append(finalArgs, patch)
	finalArgs = append(finalArgs, args...)

	sql := fmt.Sprintf(
		"UPDATE %s SET data = COALESCE(data, '{}'::jsonb) || $1::jsonb",
		quoteTable(table),
	)

	if whereClause != "" {
		sql += " WHERE " + whereClause
	}

	tag, err := e.runner.Exec(ctx, sql, finalArgs...)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func (e *pgxExecutor) Delete(ctx context.Context, table string, where *builder.WhereBuilder[json.RawMessage]) (int64, error) {
	whereClause, args, err := e.buildWhereFromWhere(where, 1)
	if err != nil {
		return 0, err
	}

	sql := fmt.Sprintf("DELETE FROM %s", quoteTable(table))
	if whereClause != "" {
		sql += " WHERE " + whereClause
	}

	tag, err := e.runner.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

// #endregion
// #region transação (PgxTransaction)

type PgxTransaction struct {
	tx pgx.Tx
}

var _ common.IUnitOfWork = (*PgxTransaction)(nil)

func (t *PgxTransaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *PgxTransaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

// #endregion
// #region PgxDatabaseAdapter (IDatabaseAdapter)

type PgxDatabaseAdapter struct {
	pool *pgxpool.Pool
	exec *pgxExecutor
}

var _ database.IDatabaseAdapter = (*PgxDatabaseAdapter)(nil)

func NewPgxDatabaseAdapter(databaseConfig *database.DatabaseConfig) *PgxDatabaseAdapter {
	if databaseConfig == nil {
		panic("NewPgxDatabaseAdapter: databaseConfig is nil")
	}

	pgxConfig, err := pgxpool.ParseConfig(databaseConfig.PostgresURI)
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		panic(err)
	}

	return &PgxDatabaseAdapter{
		pool: pool,
		exec: &pgxExecutor{
			config: databaseConfig,
			runner: pool,
		},
	}
}

func (p *PgxDatabaseAdapter) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

func (p *PgxDatabaseAdapter) Ping(ctx context.Context) error {
	if p.pool == nil {
		return errors.New("PgxDatabaseAdapter: pool is nil")
	}
	return p.pool.Ping(ctx)
}

func (p *PgxDatabaseAdapter) BeginTransaction(ctx context.Context) (common.IUnitOfWork, error) {
	if p.pool == nil {
		return nil, errors.New("PgxDatabaseAdapter: pool is nil")
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &PgxTransaction{tx: tx}, nil
}

func (p *PgxDatabaseAdapter) pickExecutor(optionalUow ...common.IUnitOfWork) (*pgxExecutor, error) {
	if len(optionalUow) == 0 || optionalUow[0] == nil {
		return p.exec, nil
	}

	tx, ok := optionalUow[0].(*PgxTransaction)
	if !ok {
		return nil, fmt.Errorf("PgxDatabaseAdapter: unexpected transaction type %T", optionalUow[0])
	}

	return &pgxExecutor{
		config: p.exec.config,
		runner: tx.tx,
	}, nil
}

func (p *PgxDatabaseAdapter) FindOne(ctx context.Context, table string, query *builder.Query[json.RawMessage], optionalUow ...common.IUnitOfWork) (*json.RawMessage, error) {
	exec, err := p.pickExecutor(optionalUow...)
	if err != nil {
		return nil, err
	}
	return exec.FindOne(ctx, table, query)
}

func (p *PgxDatabaseAdapter) FindMany(ctx context.Context, table string, query *builder.Query[json.RawMessage], optionalUow ...common.IUnitOfWork) (*builder.Result[json.RawMessage], error) {
	exec, err := p.pickExecutor(optionalUow...)
	if err != nil {
		return nil, err
	}
	return exec.FindMany(ctx, table, query)
}

func (p *PgxDatabaseAdapter) Count(ctx context.Context, table string, query *builder.Query[json.RawMessage], optionalUow ...common.IUnitOfWork) (int64, error) {
	exec, err := p.pickExecutor(optionalUow...)
	if err != nil {
		return 0, err
	}
	return exec.Count(ctx, table, query)
}

func (p *PgxDatabaseAdapter) Insert(ctx context.Context, table string, entities []json.RawMessage, optionalUow ...common.IUnitOfWork) error {
	exec, err := p.pickExecutor(optionalUow...)
	if err != nil {
		return err
	}
	return exec.Insert(ctx, table, entities)
}

func (p *PgxDatabaseAdapter) Update(ctx context.Context, table string, where *builder.WhereBuilder[json.RawMessage], update *builder.UpdateBuilder[json.RawMessage], optionalUow ...common.IUnitOfWork) (int64, error) {
	exec, err := p.pickExecutor(optionalUow...)
	if err != nil {
		return 0, err
	}
	return exec.Update(ctx, table, where, update)
}

func (p *PgxDatabaseAdapter) Delete(ctx context.Context, table string, where *builder.WhereBuilder[json.RawMessage], optionalUow ...common.IUnitOfWork) (int64, error) {
	exec, err := p.pickExecutor(optionalUow...)
	if err != nil {
		return 0, err
	}
	return exec.Delete(ctx, table, where)
}

// #endregion
