HELLO, WELCOME TO THE 'gator project'

Objectives:
    This project is a CLI to 'link to' RSS feeds to have constant updates of your favorite websites or any site using RSS feeds and simply display them in your console with some little 'configurations' you can set in the command.
        - It gives you a long running command (which consistently fetches data from RSS and have a customizable refresh rate), so that you are never late on your favorite content sites. 
        - It saves its configuration file on a "~/.gatorconfig.json" file
        - It support 'multi-users' (in the same local computer)
        - Each user can follow different feed depending on their interests

For this project to run, you will need **postgres** installed in your pc.

The project uses:
    - goose (for database migrations)
    - sqlc (for code generation based on SQL commands)

Installing the required dependencies:  $go get .

Build and run:
    - $ go run build
    - $ ./blog_aggregator register "My Username"
    - $ go install
A few commands:
    - register <name>: will register a new user and set him as the current logged in user.
    - login <name>: will change the logged in user to the specified user
    - addfeed <name> <url>: will register the feed (saving it in DB, and make the current user follow it).
    - follow <url>: will make the current user follow a specific feed (if it isn't already the case).
    - agg <duration>: will constently run to save the feeds' post in our DB.
    - feeds: will show the list of feeds in the Database.
    - following: will show the list of feeds the current user is following.
    - browse: will show the list of posts related to the feeds the current user is following.

