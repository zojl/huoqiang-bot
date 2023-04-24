package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "strings"

    "jiqiren/bot/messageHandler"
    "jiqiren/bot/database"

    "github.com/SevereCloud/vksdk/v2/api"
    "github.com/SevereCloud/vksdk/v2/events"
    "github.com/SevereCloud/vksdk/v2/callback"
)

func main() {
    token := os.Getenv("TOKEN")
    prefix := os.Getenv("PREFIX")
    vk := api.NewVK(token)

    database.Init()

    // get information about the group
    _, err := vk.GroupsGetByID(nil)
    if err != nil {
        log.Fatal(err)
    }

    cb := callback.NewCallback()
    cb.ConfirmationKey = os.Getenv("CALLBACK_RESPONSE")
    cb.SecretKey = os.Getenv("CALLBACK_SECRET")

    cb.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
        if (os.Getenv("ENV") == "dev") {
            log.Printf("%+v\n", obj)
        }

        if (strings.HasPrefix(obj.Message.Text, prefix)) {
            messageHandler.HandleCommand(&obj, vk)
        } else {
            messageHandler.HandlePlain(obj, vk)
        }
    })

    log.Println("Starting Web Server")
    http.HandleFunc(os.Getenv("CALLBACK_URL"), cb.HandleFunc)
    if err := http.ListenAndServe(":" + os.Getenv("CALLBACK_PORT"), nil); err != nil {
        log.Fatal(err)
    }
}
