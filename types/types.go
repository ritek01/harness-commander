package types

type EntityType string
type ImportType string
type StoreType int64

const (
	UNDEFINED StoreType = iota
	INLINE
	REMOTE
)

type Filter struct {
	Type           ImportType `json:"importType"`
	AppId          string     `json:"appId"`
	TriggerIds     []string   `json:"triggerIds"`
	WorkflowIds    []string   `json:"workflowIds"`
	PipelineIds    []string   `json:"pipelineIds"`
	Ids            []string   `json:"ids"`
	ServiceIds     []string   `json:"serviceIds"`
	EnvironmentIds []string   `json:"environmentIds"`
}

type DestinationDetails struct {
	AccountIdentifier string `json:"accountIdentifier"`
	AuthToken         string `json:"authToken"`
	ProjectIdentifier string `json:"projectIdentifier"`
	OrgIdentifier     string `json:"orgIdentifier"`
}

type EntityDefaults struct {
	Scope              string `json:"scope"`
	WorkflowAsPipeline bool   `json:"workflowAsPipeline"`
}

type Defaults struct {
	SecretManagerTemplate EntityDefaults `json:"SECRET_MANAGER_TEMPLATE"`
	SecretManager         EntityDefaults `json:"SECRET_MANAGER"`
	Secret                EntityDefaults `json:"SECRET"`
	Connector             EntityDefaults `json:"CONNECTOR"`
	Workflow              EntityDefaults `json:"WORKFLOW"`
	Template              EntityDefaults `json:"TEMPLATE"`
}

type Inputs struct {
	Defaults Defaults `json:"defaults"`
}

type OrgDetails struct {
	Identifier  string `json:"identifier"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectDetails struct {
	OrgIdentifier string   `json:"orgIdentifier"`
	Identifier    string   `json:"identifier"`
	Name          string   `json:"name"`
	Color         string   `json:"color"`
	Modules       []string `json:"modules"`
	Description   string   `json:"description"`
}

type ProjectBody struct {
	Project ProjectDetails `json:"project"`
}

type ProjectListBody struct {
	Projects []ProjectBody `json:"content"`
}

type OrgListBody struct {
	Organisations []OrgResponse `json:"content"`
}

type OrgResponse struct {
	Org OrgBody `json:"organizationResponse"`
}

type OrgBody struct {
	Org OrgDetails `json:"organization"`
}

type RequestBody struct {
	DestinationDetails   DestinationDetails `json:"destinationDetails"`
	EntityType           EntityType         `json:"entityType"`
	Filter               Filter             `json:"filter"`
	Inputs               Inputs             `json:"inputs"`
	IdentifierCaseFormat string             `json:"identifierCaseFormat"`
}

type CurrentGenEntity struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	AppId string `json:"appId"`
}

type SkipDetail struct {
	Entity CurrentGenEntity `json:"cgBasicInfo"`
	Reason string           `json:"reason"`
}

type UpgradeError struct {
	Message string           `json:"message"`
	Entity  CurrentGenEntity `json:"entity"`
}

type MigrationStats struct {
	SuccessfullyMigrated int64 `json:"successfullyMigrated"`
	AlreadyMigrated      int64 `json:"alreadyMigrated"`
}

type ResponseBody struct {
	Code     interface{}        `json:"code"`
	Message  string             `json:"message"`
	Status   string             `json:"status"`
	Data     interface{}        `json:"data"`
	Resource interface{}        `json:"resource"`
	Messages []ResponseMessages `json:"responseMessages"`
}

type ResponseMessages struct {
	Code         string      `json:"code"`
	Level        string      `json:"level"`
	Message      string      `json:"message"`
	Exception    interface{} `json:"exception"`
	FailureTypes interface{} `json:"failureTypes"`
}

type SummaryDetails struct {
	Count  int64  `json:"count"`
	Status string `json:"status"`
}

type SecretStore struct {
	ApiKey    string `json:"apiKey"`
	AccountId string `json:"accountId"`
	BaseURL   string `json:"baseUrl"`
	UserId    string `json:"userId"`
}
type SecretSpec struct {
	Value                   string `json:"value,omitempty"`
	SecretManagerIdentifier string `json:"secretManagerIdentifier"`
	ValueType               string `json:"valueType"`
}
type Secret struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
	Description string `json:"description,omitempty"`
	Tags        struct {
	} `json:"tags,omitempty"`
	OrgIdentifier     string      `json:"orgIdentifier,omitempty"`
	ProjectIdentifier string      `json:"projectIdentifier,omitempty"`
	Spec              interface{} `json:"spec"`
}

type HarnessSecret struct {
	Secret `json:"secret"`
}

type HarnessService struct {
	Identifier        string `json:"identifier"`
	Name              string `json:"name"`
	ProjectIdentifier string `json:"projectIdentifier,omitempty"`
	OrgIdentifier     string `json:"orgIdentifier,omitempty"`
	Description       string `json:"description,omitempty"`
	Tags              struct {
	} `json:"tags,omitempty"`
	Yaml string `json:"yaml"`
}

type HarnessEnvironment struct {
	Identifier        string `json:"identifier"`
	Name              string `json:"name"`
	ProjectIdentifier string `json:"projectIdentifier,omitempty"`
	OrgIdentifier     string `json:"orgIdentifier,omitempty"`
	Description       string `json:"description,omitempty"`
	Tags              struct {
	} `json:"tags,omitempty"`
	Color string `json:"color,omitempty"`
	Type  string `json:"type,omitempty"`
	Yaml  string `json:"yaml"`
}

type HarnessInfra struct {
	Identifier        string `json:"identifier"`
	Name              string `json:"name"`
	ProjectIdentifier string `json:"projectIdentifier,omitempty"`
	OrgIdentifier     string `json:"orgIdentifier,omitempty"`
	Description       string `json:"description,omitempty"`
	Tags              struct {
	} `json:"tags,omitempty"`
	Yaml string `json:"yaml"`
}

type HarnessPipeline struct {
	Identifier        string `json:"identifier"`
	Name              string `json:"name"`
	ProjectIdentifier string `json:"projectIdentifier,omitempty"`
	OrgIdentifier     string `json:"orgIdentifier,omitempty"`
	Description       string `json:"description,omitempty"`
	Tags              struct {
	} `json:"tags,omitempty"`
	RepoIdentifier string    `json:"repoIdentifier,omitempty"`
	RootFolder     string    `json:"rootFolder,omitempty"`
	FilePath       string    `json:"filePath,omitempty"`
	BaseBranch     string    `json:"baseBranch,omitempty"`
	ConnectorRef   string    `json:"connectorRef,omitempty"`
	RepoName       string    `json:"repoName,omitempty"`
	StoreType      StoreType `json:"storeType"`
	Yaml           string    `json:"yaml"`
}

const (
	SecretText string = "SecretText"
	SecretFile        = "SecretFile"
	SSHKey            = "SSHKey"
	WinRM             = "WinRmCredentials"
)

type CliRequest struct {
	AuthToken   string `survey:"authToken"`
	AuthType    string `survey:"authType"`
	Account     string `survey:"account"`
	OrgName     string `survey:"default"`
	ProjectName string `survey:"default"`
	Debug       bool   `survey:"debug"`
	Json        bool   `survey:"json"`
	BaseUrl     string `survey:"https://app.harness.io/"` //TODO : make it environment specific in utils
	UserId      string `survey:"userId"`
}

type UserInfo struct {
	Accounts []struct {
		AccountName       string `json:"accountName"`
		CompanyName       string `json:"companyName"`
		CreatedFromNG     bool   `json:"createdFromNG"`
		DefaultExperience string `json:"defaultExperience"`
		NextGenEnabled    bool   `json:"nextGenEnabled"`
		Uuid              string `json:"uuid"`
	} `json:"accounts"`
	Admin                          bool        `json:"admin"`
	BillingFrequency               interface{} `json:"billingFrequency"`
	CreatedAt                      int64       `json:"createdAt"`
	DefaultAccountId               string      `json:"defaultAccountId"`
	Disabled                       bool        `json:"disabled"`
	Edition                        interface{} `json:"edition"`
	Email                          string      `json:"email"`
	EmailVerified                  bool        `json:"emailVerified"`
	ExternalId                     interface{} `json:"externalId"`
	ExternallyManaged              bool        `json:"externallyManaged"`
	FamilyName                     interface{} `json:"familyName"`
	GivenName                      interface{} `json:"givenName"`
	Intent                         interface{} `json:"intent"`
	LastUpdatedAt                  int64       `json:"lastUpdatedAt"`
	Locked                         bool        `json:"locked"`
	Name                           string      `json:"name"`
	SignupAction                   interface{} `json:"signupAction"`
	Token                          interface{} `json:"token"`
	TwoFactorAuthenticationEnabled bool        `json:"twoFactorAuthenticationEnabled"`
	UtmInfo                        struct {
		UtmCampaign string `json:"utmCampaign"`
		UtmContent  string `json:"utmContent"`
		UtmMedium   string `json:"utmMedium"`
		UtmSource   string `json:"utmSource"`
		UtmTerm     string `json:"utmTerm"`
	} `json:"utmInfo"`
	Uuid string `json:"uuid"`
}

type Connector struct {
	AccountIdentifier string `json:"accountIdentifier"`
	Identifier        string `json:"identifier"`
	Name              string `json:"name"`
	ProjectIdentifier string `json:"projectIdentifier,omitempty"`
	OrgIdentifier     string `json:"orgIdentifier,omitempty"`
	Description       string `json:"description,omitempty"`
	ConnType          string `json:"type"`
}
type ConnectorListBody struct {
	Connectors []ConnectorResponse `json:"content"`
}

type ConnectorResponse struct {
	Conn ConnectorBody `json:"connectorResponse"`
}

type ConnectorBody struct {
	Conn Connector `json:"connector"`
}

type CreateRepoRequest struct {
	DefaultBranch string `json:"default_branch"`
	Description   string `json:"description"`
	GitIgnore     string `json:"git_ignore"`
	IsPublic      bool   `json:"is_public"`
	License       string `json:"license"`
	Uid           string `json:"uid"`
	Readme        bool   `json:"readme"`
	ParentRef     string `json:"parent_ref"`
}

type CreateRepoResponse struct {
	ID             int    `json:"id"`
	ParentID       int    `json:"parent_id"`
	Identifier     string `json:"identifier"`
	Path           string `json:"path"`
	Description    string `json:"description"`
	IsPublic       bool   `json:"is_public"`
	CreatedBy      int    `json:"created_by"`
	Created        int64  `json:"created"`
	Updated        int64  `json:"updated"`
	Size           int    `json:"size"`
	SizeUpdated    int    `json:"size_updated"`
	DefaultBranch  string `json:"default_branch"`
	ForkID         int    `json:"fork_id"`
	NumForks       int    `json:"num_forks"`
	NumPulls       int    `json:"num_pulls"`
	NumClosedPulls int    `json:"num_closed_pulls"`
	NumOpenPulls   int    `json:"num_open_pulls"`
	NumMergedPulls int    `json:"num_merged_pulls"`
	Importing      bool   `json:"importing"`
	IsEmpty        bool   `json:"is_empty"`
	GitURL         string `json:"git_url"`
	UID            string `json:"uid"`
}
