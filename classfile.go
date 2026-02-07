package gitrans

const (
	XGoPackage = true
)

// -----------------------------------------------------------------------------

type handler struct {
	pattern  []*patternImpl
	callback func(f *File)
}

// App represents the main application.
type App struct {
	handlers []handler
	upstream string // upstream branch
	noExec   bool   // dry run
}

// iAppProto represents the interface for application prototype.
type iAppProto interface {
	initApp() *App
	MainEntry()
}

func (p *App) initApp() *App {
	p.upstream = "main"
	return p
}

// From sets the upstream branch.
func (p *App) From(upstreamBranch string) {
	p.upstream = upstreamBranch
}

// DryRun sets the dry run mode.
func (p *App) DryRun(b bool) {
	p.noExec = b
}

// OnFile registers a callback to be executed on each file matching the pattern.
func (p *App) OnFile__0(pattern string, callback func(f *File)) {
	p.handlers = append(p.handlers, handler{
		pattern:  []*patternImpl{parsePattern(pattern)},
		callback: callback,
	})
}

// OnFile registers a callback to be executed on each file matching any of the patterns.
func (p *App) OnFile__1(patterns []string, callback func(f *File)) {
	p.handlers = append(p.handlers, handler{
		pattern:  parsePatterns(patterns),
		callback: callback,
	})
}

// XGot_App_Main is required by XGo compiler as the entry of a git_patch project.
func XGot_App_Main(app iAppProto) {
	theApp := app.initApp()
	app.MainEntry()
	newApplyer(theApp).run()
}

// -----------------------------------------------------------------------------
