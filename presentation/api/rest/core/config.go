package core

import "src/core/validator"

type Config struct {
	Port                int    `json:"port"`      // ex: 3000
	BasePath            string `json:"base_path"` // ex: "/api"
	EnableSwagger       bool   `json:"enable_swagger"`
	SwaggerPath         string `json:"swagger_path"` // ex: "/docs"
	SwaggerTitle        string `json:"swagger_title"`
	SwaggerDescription  string `json:"swagger_description"`
	SwaggerContactName  string `json:"swagger_contact_name"`
	SwaggerContactURL   string `json:"swagger_contact_url"`
	SwaggerContactEmail string `json:"swagger_contact_email"`
	SwaggerLicenseName  string `json:"swagger_license_name"`
	SwaggerLicenseURL   string `json:"swagger_license_url"`
	SwaggerVersion      string `json:"swagger_version"`
}

var _ validator.IValidable = (*Config)(nil)

func (config *Config) Validate() error {
	return validator.Object(config,
		validator.Number(&config.Port).Required().Integer().Positive().Default(4000),
		validator.String(&config.BasePath).Required().Default(""),
		validator.Boolean(&config.EnableSwagger).Default(true),
		validator.String(&config.SwaggerPath).Required().Default("/"),
		validator.String(&config.SwaggerTitle).Required().Default("Bflow - Control Plane API"),
		validator.String(&config.SwaggerDescription).Required().Default("API for Control Plane of Bflow solution"),
		validator.String(&config.SwaggerContactName).Required().Default("Leandro Santiago Gomes"),
		validator.String(&config.SwaggerContactURL).Required().Default("https://github.com/leandroluk"),
		validator.String(&config.SwaggerContactEmail).Required().Default("leandroluk@gmail.com"),
		validator.String(&config.SwaggerLicenseName).Required().Default("MIT"),
		validator.String(&config.SwaggerLicenseURL).Required().Default("https://opensource.org/licenses/MIT"),
		validator.String(&config.SwaggerVersion).Required().Default("1.0.0"),
	).Validate()
}
