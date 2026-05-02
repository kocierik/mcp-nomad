package prompts

import (
	"github.com/kocierik/mcp-nomad/utils"
)

// promptArgumentsToAny converts MCP prompt string arguments to the shape used by utils.EffectiveToolNamespace.
func promptArgumentsToAny(args map[string]string) map[string]interface{} {
	if len(args) == 0 {
		return nil
	}
	out := make(map[string]interface{}, len(args))
	for k, v := range args {
		out[k] = v
	}
	return out
}

func effectiveNamespaceFromPrompt(args map[string]string) string {
	return utils.EffectiveToolNamespace(promptArgumentsToAny(args))
}
