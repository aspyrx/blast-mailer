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

###Recipients

The recipients file should be a comma-separated plain text file.  
The first line in the file should be a list of tags that define the keywords that will be substituted later for actual values.  
The rows that follow contain the values that will be substituted. They should be in the same order as the tags - i.e., the first tag corresponds to the first value, the second tag to the second value, etc.

**Note**: The first value in each row _must_ be the recipient's email.

For example, a complete recipients file:

    EMAIL, FIRST_NAME, LAST_NAME, UNIQUE_MESSAGE
    alice@example.com, Alice, Smith, foo
    john@example.com, John, Smith, bar
    
###Message

The message file should be a plaintext SMTP request body, headers optional. Any two "$" characters surrounding a tag specified in the recipients file will be replaced with the appropriate value for that field. For example, using the previous recipients file, the message:

    Content-Type: text/html
    To: "$FIRST_NAME$ $LAST_NAME" <$EMAIL$>
    Subject: blast-mailer Test
    
    <p>Hi, $FIRST_NAME$!</p>
    <p>Your email address is $EMAIL$.</p>
    <p>Here is a unique message for you: $UNIQUE_MESSAGE$</p>
    
would be sent to the aforementioned Alice Smith as follows:

    Content-Type: text/html
    To: "Alice Smith" <alice@example.com>
    Subject: blast-mailer Test
    
    <p>Hi, Alice!</p>
    <p>Your email address is alice@example.com.</p>
    <p>Here is a unique message for you: foo</p>
    
