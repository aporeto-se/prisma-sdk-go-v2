package common

// PrismaToken is a JWT (jwt.io) token issued by Prisma for access to the Prisma API
type PrismaToken struct {
	Audience string `json:"audience,omitempty" yaml:"audience,omitempty"`
	Claims   struct {
		Data struct {
			Account         string `json:"account,omitempty" yaml:"account,omitempty"`
			Arn             string `json:"arn,omitempty" yaml:"arn,omitempty"`
			Organization    string `json:"organization,omitempty" yaml:"organization,omitempty"`
			Partition       string `json:"partition,omitempty" yaml:"partition,omitempty"`
			Realm           string `json:"realm,omitempty" yaml:"realm,omitempty"`
			Resource        string `json:"resource,omitempty" yaml:"resource,omitempty"`
			Resourcetype    string `json:"resourcetype,omitempty" yaml:"resourcetype,omitempty"`
			Rolename        string `json:"rolename,omitempty" yaml:"rolename,omitempty"`
			Rolesessionname string `json:"rolesessionname,omitempty" yaml:"rolesessionname,omitempty"`
			Service         string `json:"service,omitempty" yaml:"service,omitempty"`
			Subject         string `json:"subject,omitempty" yaml:"subject,omitempty"`
			Userid          string `json:"userid,omitempty" yaml:"userid,omitempty"`
			Email           string `json:"email,omitempty" yaml:"email,omitempty"`
			Instancename    string `json:"instancename,omitempty" yaml:"instancename,omitempty"`
			Projectid       string `json:"projectid,omitempty" yaml:"projectid,omitempty"`
			Projectnumber   string `json:"projectnumber,omitempty" yaml:"projectnumber,omitempty"`
			Zone            string `json:"zone,omitempty" yaml:"zone,omitempty"`
		} `json:"data,omitempty" yaml:"data,omitempty"`
		Exp          int64  `json:"exp,omitempty" yaml:"exp,omitempty"`
		Iat          int64  `json:"iat,omitempty" yaml:"iat,omitempty"`
		Iss          string `json:"iss,omitempty" yaml:"iss,omitempty"`
		Realm        string `json:"realm,omitempty" yaml:"realm,omitempty"`
		Restrictions struct {
		} `json:"restrictions,omitempty" yaml:"restrictions,omitempty"`
		Sub string `json:"sub,omitempty" yaml:"sub,omitempty"`
	} `json:"claims,omitempty" yaml:"claims,omitempty"`
	Data     string      `json:"data,omitempty" yaml:"data,omitempty"`
	Metadata interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Opaque   struct {
	} `json:"opaque,omitempty" yaml:"opaque,omitempty"`
	Quota                 int           `json:"quota,omitempty" yaml:"quota,omitempty"`
	Realm                 string        `json:"realm,omitempty" yaml:"realm,omitempty"`
	RestrictedNamespace   string        `json:"restrictedNamespace,omitempty" yaml:"restrictedNamespace,omitempty"`
	RestrictedNetworks    []interface{} `json:"restrictedNetworks,omitempty" yaml:"restrictedNetworks,omitempty"`
	RestrictedPermissions []interface{} `json:"restrictedPermissions,omitempty" yaml:"restrictedPermissions,omitempty"`
	Token                 string        `json:"token,omitempty" yaml:"token,omitempty"`
	Validity              string        `json:"validity,omitempty" yaml:"validity,omitempty"`
}

// OAuthToken OAuth Token
type OAuthToken struct {
	Realm string `json:"realm"`
	Data  struct {
		CommonName   string `json:"commonName"`
		Organization string `json:"organization"`
		Realm        string `json:"realm"`
		SerialNumber string `json:"serialNumber"`
		Subject      string `json:"subject"`
	} `json:"data"`
	Restrictions struct {
	} `json:"restrictions"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
	Iss string `json:"iss"`
	Sub string `json:"sub"`
}
