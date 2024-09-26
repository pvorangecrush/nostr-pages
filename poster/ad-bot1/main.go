package main

import (
    "time"
    "fmt"
    "context"
    "os/exec"
    "os"
    "log"
    "encoding/json"
    "github.com/nbd-wtf/go-nostr"
    "github.com/nbd-wtf/go-nostr/nip19"
)

type ProfileMetadata struct {
    Name        string `json:"name"`
    DisplayName string `json:"display_name"`
    About       string `json:"about"`
    Picture     string `json:"picture"`
    Nip05       string `json:"nip05,omitempty"`
    Lud16       string `json:"lud16,omitempty"`
    Website     string `json:"website,omitempty"`
    Banner      string `json:"banner,omitempty"`
    Bot         string `json:"bot"`
}

func main() {
    sk := nostr.GeneratePrivateKey()
    pk, _ := nostr.GetPublicKey(sk)
    nsec, _ := nip19.EncodePrivateKey(sk)
    npub, _ := nip19.EncodePublicKey(pk)
    new_id := "relay-" + string(npub)[len(npub)-5:]

    directory := "/home/user/build/nostr-pages"
    file := directory + "/.well-known/nostr.json"
    file2 := directory + "/poster/posts.json"

    cmd :=exec.Command("git", "pull")
    cmd.Dir = directory
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        log.Fatalf("cmd.Run () failed with %s\n", err)
    }
    time.Sleep(10 * time.Second)

    // Add new identity to nostr.json
    sed_cmd := fmt.Sprintf(`/names/c\  "names": {\n    "%s":"%s",`, new_id, pk)
    fmt.Println("sed", "-i", sed_cmd, file)
    cmd =exec.Command("sed", "-i", sed_cmd, file)
    cmd.Dir = directory
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err = cmd.Run()
    if err != nil {
        log.Fatalf("cmd.Run () failed with %s\n", err)
    }

    // Add nostr.json to commit
    cmd =exec.Command("git", "add", "-A")
    cmd.Dir = directory
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err = cmd.Run()
    if err != nil {
        log.Fatalf("cmd.Run () failed with %s\n", err)
    }

    // Commit update to nostr.json
    cmd =exec.Command("git", "commit", "-m", "updated nostr.json")
    cmd.Dir = directory
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err = cmd.Run()
    if err != nil {
        log.Fatalf("cmd.Run () failed with %s\n", err)
    }

    // Push update to remote repository so that github-pages will publish
    cmd =exec.Command("git", "push")
    cmd.Dir = directory
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err = cmd.Run()
    if err != nil {
        log.Fatalf("cmd.Run () failed with %s\n", err)
    }
    time.Sleep(5 * time.Second)

    // Print keys for reference
    fmt.Println("sk:", sk)
    fmt.Println("pk:", pk)
    fmt.Println(nsec)
    fmt.Println(npub)
    fmt.Println(new_id)

    // Prep a post from posts.json
    newPost := advertiseRelay(file2)
    fmt.Println(newPost)

    // Prep TextNote ad event
    ev := nostr.Event{
        PubKey:    pk,
        CreatedAt: nostr.Now(),
        Kind:      nostr.KindTextNote,
        Tags:      nil,
        Content:   newPost,
    }

    // Prep Metadata Event object
    var profileMetadata ProfileMetadata
    profileMetadata.Name = new_id
    profileMetadata.DisplayName = "Orange Crush Relay Ad Bot " + new_id
    profileMetadata.About = "https://relay.orange-crush.com is a nostr relay for you! One time joining fee of 8,000 sats. Please consider us for your nostr relay needs and help keep the nostr a robust, censorship resistant protocol for everyone."
    profileMetadata.Picture = "https://pfp.nostr.build/4d6f51938709eb20733f2d0931a334f7f36446066dbd4c29d025f06107a5e1d9.png"
    profileMetadata.Nip05 = new_id + "@orange-crush.com"
    profileMetadata.Lud16 = "relay@orange-crush.com"
    profileMetadata.Website = "https://relay.orange-crush.com"
    profileMetadata.Banner = "https://m.primal.net/Hlus.jpg"
    profileMetadata.Bot = "true"
    var newContent string
    marshalledProfile, err := json.Marshal(profileMetadata)
    if err != nil {
        fmt.Println(err)
    } else {
        newContent = string(marshalledProfile)
        fmt.Println(err)
    }

    // Prep Metadata Event
    meta_ev := nostr.Event{
        PubKey:     pk,
        CreatedAt:  nostr.Now(),
        Kind:       nostr.KindProfileMetadata,
        Tags:       nil,
        Content:    newContent,
    }

    // calling Sign sets the event ID field and the event Sig field
    ev.Sign(sk)
    meta_ev.Sign(sk)

    // publish the event to relays
    ctx := context.Background()
    for _, url := range []string{"wss://relay.nostr.band", "wss://relay.mutinywallet.com", "wss://purplerelay.com", "wss://relay.damus.io", "wss://relay.snort.social"} {
        relay, err := nostr.RelayConnect(ctx, url)
        if err != nil {
            fmt.Println(err)
            continue
        }
        if err := relay.Publish(ctx, meta_ev); err != nil {
            fmt.Println(err)
            continue
        }

        fmt.Printf("published metadata note to %s\n", url)
    }
    for _, url := range []string{"wss://relay.nostr.band", "wss://relay.mutinywallet.com", "wss://purplerelay.com", "wss://relay.damus.io", "wss://relay.snort.social"} {
        relay, err := nostr.RelayConnect(ctx, url)
        if err != nil {
            fmt.Println(err)
            continue
        }
            if err := relay.Publish(ctx, ev); err != nil {
            fmt.Println(err)
            continue
        }

        fmt.Printf("published text note to %s\n", url)
    }
}