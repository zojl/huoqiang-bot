package main

import (
    "context"
    //"fmt"
    "log"
    "os"
    "strings"

    "huoqiang/bot/messageHandler"

    "github.com/SevereCloud/vksdk/v2/api"
    //"github.com/SevereCloud/vksdk/v2/api/params"
    "github.com/SevereCloud/vksdk/v2/events"
    "github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

func main() {
  token := os.Getenv("TOKEN")
  prefix := os.Getenv("PREFIX")
  vk := api.NewVK(token)

  // get information about the group
  group, err := vk.GroupsGetByID(nil)
  if err != nil {
	  log.Fatal(err)
  }

  // Initializing Long Poll
  lp, err := longpoll.NewLongPoll(vk, group[0].ID)
  if err != nil {
	  log.Fatal(err)
  }

  lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
    log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)
    if (strings.HasPrefix(obj.Message.Text, prefix)) {
      log.Printf("Received command: %s", obj.Message);
    } else {
      messageHandler.HandlePlain(obj)
    }

	  /*
        if obj.Message.Text == "ping" {
		  b := params.NewMessagesSendBuilder()
		  b.Message("pong")
		  b.RandomID(0)
		  b.PeerID(obj.Message.PeerID)

		  _, err := vk.MessagesSend(b.Params)
		  if err != nil {
			  log.Fatal(err)
		  }
	  }
      */
  })

  log.Println("Start Long Poll")
  if err := lp.Run(); err != nil {
	  log.Fatal(err)
  }
}
