package agents

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/prathyushnallamothu/swarmgo"
	"github.com/prathyushnallamothu/swarmgo/llm"
)

type Agents struct {
	Name          string
	Model         string
	Instructions  string
	provider      llm.LLMProvider
	Agent         *swarmgo.Swarm
	modelOverride string
	Stream        bool
	Debug         bool
	MaxTurns      int
	executeTools  bool
	Functions     []swarmgo.AgentFunction
}

func (a *Agents) LoadAgent() *swarmgo.Agent {
	agent := &swarmgo.Agent{
		Name:         a.Name,
		Model:        a.Model,
		Instructions: a.Instructions,
		Provider:     a.provider,
	}
	a.initSwarm()
	return agent
}

func (a *Agents) initSwarm() {
	a.Agent = swarmgo.NewSwarm(a.loadApiKey(), a.provider)
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

	if apiKey, ok := apiKeyMap[a.provider]; ok {
		return apiKey
	}
	return ""

}
