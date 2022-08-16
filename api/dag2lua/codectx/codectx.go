package codectx

import (
	"fmt"
	"strings"

	"github.com/apisix/manager-api/dag2lua/utils"
)

type CodeCtx struct {
	globals    map[string]string
	preface    []string
	params     []string
	nLabels    int
	userValues []string
	body       []string
}

func (ctx *CodeCtx) Libfunc(globalName string) string {
	localName, ok := ctx.globals[globalName]

	if !ok {
		localName = strings.Replace(globalName, "%.", "_", 1)
		ctx.globals[globalName] = localName
		ctx.AddPreface(fmt.Sprintf("local %v = %v", localName, globalName))
	}
	return localName
}

func (ctx *CodeCtx) AddParam(param string) string {
	ctx.params = append(ctx.params, param)
	return param
}

func (ctx *CodeCtx) AddLabel() string {
	ctx.nLabels++
	return "label_" + fmt.Sprint(ctx.nLabels)
}

func (ctx *CodeCtx) UserValue(val string) string {
	slot := len(ctx.userValues) + 1
	ctx.userValues[slot] = val
	return fmt.Sprintf("uservalues[%d]", slot)
}

func (ctx *CodeCtx) Generate(rule, conf string) string {
	ruleCtx, err := generateRule(ctx.Child(), rule, conf)
	if err != nil {
		return nil, err
	}

	ctx.Stmt(fmt.Sprintf("%s = ", "_M.access"), ruleCtx, "\n\n")
	ctx.Stmt(fmt.Sprintf("%s = ", "_M.header_filter"),
		GenerateCommonPhase(ctx.Child(), "header_filter"), "\n\n")
	ctx.Stmt(fmt.Sprintf("%s = ", "_M.body_filter"),
		GenerateCommonPhase(ctx.Child(), "body_filter"), "\n\n")

	release_plugins := `tablepool.release("script_plugins", ctx.script_plugins)`
	ctx.Stmt(fmt.Sprintf("%s = ", "_M.log"),
		GenerateCommonPhase(ctx.Child(), "log", release_plugins), "\n\n")

	return "_M"
}

func (ctx *CodeCtx) AddPreface(args ...string) {
	for i := 0; i < len(args); i++ {
		ctx.preface = append(ctx.preface, args[i])
	}
	ctx.preface = append(ctx.preface, "\n")
}

func (ctx *CodeCtx) Stmt(args ...string) {
	for i := 0; i < len(args); i++ {
		ctx.body = append(ctx.body, args[i])
	}
	ctx.body = append(ctx.body, "\n")
}

func (ctx *CodeCtx) generate(codeTable []string) {
	indent := " "
	for _, stmt := range ctx.preface {
		utils.InsertCode(indent, codeTable)
		if getmetatable(stmt) == ctx {

		}
	}
}

func (ctx *CodeCtx) GenerateCommonPhase(phase string, tailLua string) *codectx.CodeCtx {
	ctx.Stmt("local plugins = ctx.script_plugins")
	ctx.Stmt("for i = 1, #plugins, 2 do")
	ctx.Stmt("    local plugin_name = plugins[i]")
	ctx.Stmt("    local plugin_conf_name = plugins[i + 1]")
	ctx.Stmt("    local plugin_obj = plugin.get(plugin_name)")
	ctx.Stmt("    local phase_fun = plugin_obj." + phase)
	ctx.Stmt("    if phase_fun then")
	ctx.Stmt(fmt.Sprintf("        phase_fun(_M[plugin_conf_name], %s)", ctx.AddParam("ctx")))
	ctx.Stmt("    end")
	ctx.Stmt("end")
	if tailLua != "" {
		ctx.Stmt(tailLua)
	}
	return ctx
}

// func (ctx *CodeCtx) Child(ref string) interface {
// 	return map {
// 		schema = ref

// 	}
// }
