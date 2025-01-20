package agents

import (
	"context"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/prathyushnallamothu/swarmgo"
	"github.com/prathyushnallamothu/swarmgo/llm"
)

type Agents struct {
	Name          string
	Model         string
	Instructions  string
	Provider      llm.LLMProvider
	Core          *swarmgo.Swarm
	Agent         *swarmgo.Agent
	ModelOverride string
	Stream        bool
	Debug         bool
	MaxTurns      int
	ExecuteTools  bool
	Functions     []swarmgo.AgentFunction
}

func (a *Agents) LoadAgent() {
	agent := &swarmgo.Agent{
		Name:         a.Name,
		Model:        a.Model,
		Instructions: a.Instructions,
		Provider:     a.Provider,
	}
	a.Agent = agent
	a.initSwarm()
}

func (a *Agents) initSwarm() {
	a.Core = swarmgo.NewSwarm(a.loadApiKey(), a.Provider)
}

func (a *Agents) loadApiKey() string {
	openaiApiKey, _ := beego.AppConfig.String("OPENAI_API_KEY")
	azureApiKey, _ := beego.AppConfig.String("AZURE_API_KEY")
	azureAdApiKey, _ := beego.AppConfig.String("AZURE_AD_API_KEY")
	cloudflareAzureApiKey, _ := beego.AppConfig.String("CLOUDFLARE_AZURE_API_KEY")
	geminiApiKey, _ := beego.AppConfig.String("GEMINI_API_KEY")
	claudeApiKey, _ := beego.AppConfig.String("CLAUDE_API_KEY")
	ollamaApiKey, _ := beego.AppConfig.String("OLLAMA_API_KEY")
	deepSeekApiKey, _ := beego.AppConfig.String("DEEPSEEK_API_KEY")

	apiKeyMap := map[llm.LLMProvider]string{
		llm.OpenAI:          openaiApiKey,
		llm.Azure:           azureApiKey,
		llm.AzureAD:         azureAdApiKey,
		llm.CloudflareAzure: cloudflareAzureApiKey,
		llm.Gemini:          geminiApiKey,
		llm.Claude:          claudeApiKey,
		llm.Ollama:          ollamaApiKey,
		llm.DeepSeek:        deepSeekApiKey,
	}

	if apiKey, ok := apiKeyMap[a.Provider]; ok {
		return apiKey
	}
	return ""

}

func (a *Agents) Execute(ctx context.Context, message []llm.Message, contextVariables map[string]interface{}) (swarmgo.Response, error) {
	return a.Core.Run(ctx, a.Agent, message, contextVariables, a.ModelOverride, a.Stream, a.Debug, a.MaxTurns, a.ExecuteTools)
}
