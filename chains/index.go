package chains

import (
	"context"
	"errors"
	"log"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

type Vacation struct {
	Id        uuid.UUID `json:"id"`
	Completed bool      `json:"completed"`
	Idea      string    `json:"idea"`
}

var Vacations []*Vacation

func GetVacationFromDb(id uuid.UUID) (Vacation, error) {
	// Use the slices package to find the index of the object with
	// matching ID in the database. If it does not exist, this will return
	// -1
	idx := slices.IndexFunc(Vacations, func(v *Vacation) bool { return v.Id == id })

	// handle it
	if idx < 0 {
		return Vacation{}, errors.New("ID Not Found")
	}

	// Otherwise, return the Vacation object
	return *Vacations[idx], nil
}

func GeneateVacationIdeaChange(id uuid.UUID, budget int, season string, hobbies []string) {
	log.Printf("Generating new vacation with ID: %s", id)

	v := &Vacation{Id: id, Completed: false, Idea: ""}
	Vacations = append(Vacations, v)

	ctx := context.Background()
	llm, err := openai.New()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	system_message_prompt_string := "You are an AI dating coach that will help me plan a perfect date.\n" +
		"My girlfriend's favorite season is {{.season}}.\n" +
		"My girlfriend's hobbies include {{.hobbies}}.\n" +
		"My girlfriend's budget is {{.budget}} dollars.\n" +
		"My girlfriend's name is hwangjunghye.\n" +
		"Suggest a romantic date idea in Korea that aligns with her interests and the season. use korean"

	system_message_prompt := prompts.NewSystemMessagePromptTemplate(system_message_prompt_string, []string{"season", "hobbies", "dollars"})

	// Create a human prompt with the request that a human would have
	human_message_prompt_string := "write a travel itinerary for me"
	human_message_prompt := prompts.NewHumanMessagePromptTemplate(human_message_prompt_string, []string{})

	chat_prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{system_message_prompt, human_message_prompt})

	vals := map[string]any{
		"season":  season,
		"budget":  budget,
		"hobbies": strings.Join(hobbies, ","),
	}
	msgs, err := chat_prompt.FormatMessages(vals)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	content := []llms.MessageContent{
		llms.TextParts(msgs[0].GetType(), msgs[0].GetContent()),
		llms.TextParts(msgs[1].GetType(), msgs[1].GetContent()),
	}

	// Invoke the LLM with the messages which
	completion, err := llm.GenerateContent(ctx, content)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	v.Idea = completion.Choices[0].Content
	v.Completed = true

	log.Printf("Generation for %s is done!", v.Id)
}
