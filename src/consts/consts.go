package consts

const AppBanner = `
█▄▀ █▀▄ ██▀ █▀▄ ██▀ █▄ █ █ █
█ █ █▀▄ █▄▄ █▄▀ █▄▄ █ ▀█ ▀▄▀
`
const AppName = "kredenv"
const AppDesc = "Inject env vars & secrets into your shell environment"

var AppVersion = "dev"
var BuildCommit = "none"
var BuildDate = "unknown"

const RootDirName = ".kredenv"
const AuthEnvVar = "KREDENTIAL"
const AuthMasterFile = ".kredmaster"
const KeyringKey = "kredkey"
const EncFileName = ".kreds.enc"

var SupportedExportFormats = []string{"env", "json", "yaml", "toml"}
var SupportedInjectFormats = []string{"dotenv", "json"}
