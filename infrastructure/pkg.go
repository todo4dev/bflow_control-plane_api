package infra

import (
	_ "src/infrastructure/cache"
	_ "src/infrastructure/crypto"
	_ "src/infrastructure/database"
	_ "src/infrastructure/jwt"
	_ "src/infrastructure/logger"
	_ "src/infrastructure/mailer"
	_ "src/infrastructure/openid"
	_ "src/infrastructure/storage"
	_ "src/infrastructure/stream"
)
