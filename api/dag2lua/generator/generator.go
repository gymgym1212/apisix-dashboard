package generator

import (
	"errors"
)

type chart struct {
}

type conf struct {
}

type rule struct {
}

type script struct {
	Chart chart `json:"chart"`
	Conf  conf  `json:"conf"`
	Rule  rule  `json:"rule"`
}

func generate_ctx(script map[string]interface{}) (code string, err error) {
	if _, ok := script["rule"]; !ok {
		script["rule"] = interface {}
	}
	if _, ok := script["conf"]; !ok {
		script["conf"] = interface {}
	}

	ctx, err := codectx(script["rule"], script["conf"], {})
	if err != nil {
		return "0", errors.New("Get codectx error",err)
	}

	ctx.preface(`local core = require("apisix.core")`)
	ctx.preface(`local plugin = require("apisix.plugin")`)
	ctx.preface(`local tablepool = core.tablepool`)
	ctx.preface('\n')

	class_name, err := ctx:generate(data.rule, data.conf)
	if err != nil {
		return "", err
	}

	ctx.stmt(`return `, class_name)

	return ctx, nil
}
func Generate(script map[string]interface{}) (string, error) {
	ctx, err := generate_ctx(script)
	if err != nil {
		return "", err
	}
	code, err : ctx.as_lua()
	if err != nil {
		return "", err
	}
	return code, nil
}
