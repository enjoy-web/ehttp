package ehttp

type Scheme string

const (
	SchemeHTTP  Scheme = "http"
	SchemeHTTPS Scheme = "https"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
)

const (
	InHeader   = "header"
	InPath     = "path"
	InQuery    = "query"
	InFormData = "formData"
	InBody     = "body"
)

const (
	Application_ActiveMessage = "application/activemessage"
	Application_AppleFile     = "application/applefile"
	Application_AtomicMail    = "application/atomicmail"
	Application_MsWord        = "application/msword"
	Application_OctetStream   = "application/octet-stream"
	Application_Oda           = "application/oda"
	Application_Pdf           = "application/pdf"
	Application_PostScript    = "application/postscript"
	Application_Rtf           = "application/rtf"
	Application_Xml           = "application/xml"
	Application_Json          = "application/json"
	Application_Json_utf8     = "application/json;charset=utf-8"
	Application_Zip           = "application/zip"
	Application_Siren_Json    = "application/vnd.siren+json"
	Application_Hal_Json      = "application/hal+json"
	Audio_Xaiff               = "audio/x-aiff"
	Audio_Xwav                = "audio/x-wav"
	Image_Cgm                 = "image/cgm"
	Image_G3Fax               = "image/g3fax"
	Image_Gif                 = "image/gif"
	Image_Ief                 = "image/ief"
	Image_Jpeg                = "image/jpeg"
	Image_Naplps              = "image/naplps"
	Image_Png                 = "image/png"
	Image_Tiff                = "image/tiff"
	Multipart_Alternative     = "multipart/alternative"
	Multipart_AppleDouble     = "multipart/appledouble"
	Multipart_Digest          = "multipart/digest"
	Multipart_FormData        = "multipart/form-data"
	Multipart_HeaderSet       = "multipart/header-set"
	Multipart_Mixed           = "multipart/mixed"
	Multipart_Parallel        = "multipart/parallel"
	Multipart_Related         = "multipart/related"
	Multipart_Report          = "multipart/report"
	Multipart_VoiceMessage    = "multipart/voice-message"
	Text_Enriched             = "text/enriched"
	Text_Html                 = "text/html"
	Text_Plain                = "text/plain"
	Text_RichText             = "text/richtext"
	Text_Sgml                 = "text/sgml"
	Text_TabSeparatedValues   = "text/tab-separated-values"
	Text_Xml                  = "text/xml"
	Text_X_SeText             = "text/x-setext"
	Video_Mpeg                = "video/mpeg"
	Video_Quicktime           = "video/quicktime"
	Video_VndVivo             = "video/vnd.vivo"
	Video_VndMotorolaVideo    = "video/vnd.motorola.video"
	Video_VndMotorolaVideoP   = "video/vnd.motorola.videop"
	Video_X_MS_Video          = "video/x-msvideo"
	Video_X_SgiMovie          = "video/x-sgi-movie"
)
