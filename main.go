package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/nlopes/slack"
	"github.com/aws/aws-lambda-go/lambda"
)

var reactionsLookup = map[string]bool{
	"grin":                          true,
	"grinning":                      true,
	"joy":                           true,
	"joy_cat":                       true,
	"laughing":                      true,
	"rolling_on_the_floor_laughing": true,
	"simple_smile":                  true,
	"slightly_smiling_face":         true,
	"smile_cat":                     true,
	"smiley_cat":                    true,
	"smiling_face_with_smiling_eyes_and_hand_covering_mouth": true,
	"smiling_imp": true,
	"smirk":       true,
	"smirk_cat":   true,
	"sweat_smile": true,
}

const iconUrl = "http://blog.edtechie.net/wp-content/uploads/2015/07/kingofcomedy.jpg"

func Handler() {
	channelID := os.Getenv("CHANNEL_ID")
	apiKey := os.Getenv("API_KEY")

	api := slack.New(apiKey)

	var amusingUsers = map[string]int{}

	t, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	oldest := strconv.FormatInt(t.Unix(), 10)
	historyParams := slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Inclusive: true,
		Limit:     100,
		Oldest:    oldest,
	}

	history, err := api.GetConversationHistory(&historyParams)
	if err != nil {
		fmt.Errorf("Unexpected error: %s", err)
		return
	}

	for _, msg := range history.Messages {
		if msg.Reactions != nil {
			for _, reaction := range msg.Reactions {
				if _, ok := reactionsLookup[reaction.Name]; ok {
					amusingUsers[msg.User] += reaction.Count
				}
			}
		}
	}

	// find all messages with reactions that are in the list
	type userToCount struct {
		Id    string
		Count int
	}

	var ss []userToCount
	for id, count := range amusingUsers {
		ss = append(ss, userToCount{id, count})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Count > ss[j].Count
	})

	var msg string
	var winner string
	var runnerUp string

	if len(ss) > 0 {
		if user, err := api.GetUserInfo(ss[0].Id); err == nil {
			c := strconv.Itoa(ss[0].Count)
			m := "Congratulations @%s!\nToday, you were the Funniest Person in the Room.\n"
			m = m + "*%s* peers applauded your comedy stylings."
			winner = fmt.Sprintf(m, user.Name, c)
		}

		if len(ss) > 1 {
			if user, err := api.GetUserInfo(ss[1].Id); err == nil {
				c := strconv.Itoa(ss[1].Count)
				m := "... @%s - you lightly tickled %s funny-bones,"
				m = m + " _but don't give up your day job ..._\n"
				runnerUp = fmt.Sprintf(m, user.Name, c)
			}
		}
	} else {
		msg = "Oh dear ... nobody was amused today ... Must be the weekend I guess ..."
	}

	params := slack.PostMessageParameters{
		AsUser:   false,
		Username: "Rupert Pupkin Speaks",
		IconURL:  iconUrl,
		Markdown: true,
	}

	if winner != "" {
		attachment := slack.Attachment{
			Text: winner,
		}
		if runnerUp != "" {
			attachment.Fields = []slack.AttachmentField{
				{
					Title: "... and to the runner up",
					Value: runnerUp,
				},
			}
		}
		params.Attachments = []slack.Attachment{attachment}
	}
	_, _, err = api.PostMessage(channelID, msg, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

func main() {
	lambda.Start(Handler)
}
