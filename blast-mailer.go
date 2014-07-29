package main

import (
    "bufio"
    "encoding/csv"
    "flag"
    "fmt"
    "io/ioutil"
    "net/smtp"
    "os"
    "regexp"
)

func main() {
    var toFilePath string
    var msgFilePath string
    var senderEmail string
    var senderPassword string
    var smtpHost string
    var smtpPort string
    var quiet bool
    var force bool

    flag.StringVar(&toFilePath, "to", "", "Path to file containing recipients in CSV format. See http://github.com/aspyrx/blast-mailer for format specifications.")
    flag.StringVar(&msgFilePath, "msg", "", "Path to file containing the message, including headers. See http://github.com/aspyrx/blast-mailer for special replaceable identifiers.")
    flag.StringVar(&senderEmail, "email", "", "The sender's email address.")
    flag.StringVar(&senderPassword, "password", "", "The sender's password.")
    flag.StringVar(&smtpHost, "host", "", "The hostname of the SMTP server to use.")
    flag.StringVar(&smtpPort, "port", "", "The SMTP server's port (without the colon).")
    flag.BoolVar(&quiet, "quiet", false, "Quiet mode - names and email addresses will not be logged. Default: off.")
    flag.BoolVar(&force, "force", false, "Force mode - all prompts and non-fatal errors will be ignored. Default: off.")

    flag.Parse()

    msg, err := ioutil.ReadFile(msgFilePath)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    toFile, err := os.Open(toFilePath)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    r := csv.NewReader(toFile)
    r.FieldsPerRecord = 0
    r.TrimLeadingSpace = true
    
    tags, err := r.Read()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    
    to, err := r.ReadAll()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if err := toFile.Close(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if !quiet {
        fmt.Printf("Sender email: %s\n", senderEmail)
        fmt.Printf("SMTP server: %s:%s\n", smtpHost, smtpPort)
    }

    if !force {
        fmt.Printf("Are you sure you want to send an email to %d recipients? (y/N): ", len(to))
        if !isOk() {
            fmt.Println("Sending cancelled.")
            os.Exit(1)
        }
    }

    auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)
    success := 0
    for _, v := range to {
        newMsg := msg
        for i, tag := range tags {
            newMsg = regexp.MustCompile(fmt.Sprintf(`\Q$%s$\E`, tag)).ReplaceAll(newMsg, []byte(v[i]))
        }

        if !quiet {
            fmt.Println(v)
        }

        recipEmail := v[0]
        err = smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), auth, senderEmail, []string{recipEmail}, newMsg)

        if err != nil && !force {
            fmt.Printf("Error sending email to '%s'. Continue? (y/N): ", recipEmail)
            if !isOk() {
                fmt.Printf("Sending cancelled. Emails sent: %d\n", success)
                os.Exit(1)
            }
        }

        success++
    }

    fmt.Printf("Sending complete. Emails sent: %d\n", success)
    return
}

func isOk() bool {
    buf := bufio.NewReader(os.Stdin)
    in, err := buf.ReadString('\n')
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    return regexp.MustCompile(`(?i)^y|yes`).MatchString(in)
}
