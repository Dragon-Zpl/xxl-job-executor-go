package xxljob

// 响应码
const (
	SuccessCode = 200
	FailureCode = 500
)

type GlueType string

const (
	GlueType_GLUE_SHELL      GlueType = "GLUE_SHELL"
	GlueType_GLUE_PYTHON     GlueType = "GLUE_PYTHON"
	GlueType_GLUE_POWERSHELL GlueType = "GLUE_POWERSHELL"
	GlueType_GLUE_NODEJS     GlueType = "GLUE_NODEJS"
	GlueType_GLUE_PHP        GlueType = "GLUE_PHP"
	GlueType_GLUE_GO         GlueType = "BEAN"
)

func (g GlueType) String() string {
	return string(g)
}

type LogIDContextKey struct{}


func (l LogIDContextKey) String() string {
	return "timing_log_id"
}


type JobNameContextKey struct{}

func (j JobNameContextKey) String() string {
	return "timing_job_name_context_key"
}