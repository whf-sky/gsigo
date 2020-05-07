package gsigo

type CmdController struct {}

// CmdInterface  定义了一些基本的命令行处理程序操作
type CmdInterface interface {
	Init()
	Prepare()
	Execute()
	Finish()
}

func (c *CmdController) Init() {}
func (c *CmdController) Prepare() {}
func (c *CmdController) Finish() {}
