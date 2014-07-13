blast-mailer
============

A command line bulk emailer tool, written in Go.

Usage
-----

    blast-mailer <flags>
  
Flags:  
  -email: The sender's email address. (required)
  -to: Path to file containing recipients in CSV format. (required)  
  -msg: Path to file containing the message, including headers. (required)  
  -host: The hostname of the SMTP server to use. (required)  
  -password: The sender's password. (required)  
  -port: The SMTP server's port (without the colon). (required)  
  -force: Force mode - all prompts and non-fatal errors will be ignored.  
  -quiet: Quiet mode - names and email addresses will not be logged.  

For example:  

    blast-mailer -to to.csv -msg msg.txt -host smtp.example.com -port 587 -email example@example.com -password foobar

File Formats
------------
The recipients file should be a CSV file, with the recipient email address as the first field, and their name as the second. For example:

    alice@example.com, Alice Smith
    john@example.com, John Smith
    
The message file should be a plaintext SMTP request body, headers optional. Any instances of "$EMAIL$" will be replaced with the recipient's email address, and any instances of "$NAME$" will be replaced with the recipient's name. For example, the message:

    Content-Type: text/html
    To: "$NAME$" <$EMAIL$>
    Subject: blast-mailer Test
    
    <p>Hi, $NAME$!</p>
    <P>Your email address is $EMAIL$.</p>
    
would be sent to the aforementioned Alice Smith as follows:

    Content-Type: text/html
    To: "Alice Smith" <alice@example.com>
    Subject: blast-mailer Test
    
    <p>Hi, Alice Smith!</p>
    <P>Your email address is alice@example.com.</p>
    
