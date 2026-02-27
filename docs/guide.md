# sp0rkle User Guide: How to use the Bot

This guide provides an overview of the commands and features available to users of the sp0rkle IRC bot.

## 1. Factoids

Factoids allow the bot to remember and recall information.

### Adding Factoids
There are two ways to teach the bot something:
- `botNick: <key> := <value>` - Stores `<value>` for `<key>`.
- `botNick: <key> :is <value>` - Stores `<key> is <value>` for `<key>`.

### Retrieving Factoids
- Simply type the `<key>` in the channel, and the bot will respond with a random value associated with that key.
- `botNick: literal <key>` - Prints all known values for a key. If there are many, the bot will ask you to do this in a private message.

### Managing Factoids
These commands operate on the **last factoid triggered** in the channel:
- `botNick: forget that` or `delete that` - Deletes the last triggered value.
- `botNick: replace that with <new value>` - Replaces the last value.
- `botNick: that =~ /regex/replacement/` - Edit the last factoid value using regex.
- `botNick: chance of that is <N>%` - Sets the probability (0-100%) that the bot will respond to this key.
- `botNick: fact info` - Displays metadata about the last factoid (creator, time, etc.).
- `botNick: fact search <regex>` - Search for factoid keys matching a regex.

### Factoid Variables
You can use variables in factoid values:
- `$nick`: The nick of the person who triggered the factoid.
- `$chan`: The channel it was triggered in.
- `$date` / `$time`: Current date or time.
- `$user` / `$host`: Ident or host of the triggerer.

## 2. Karma

Karma tracks the "score" of things or people.

- `<thing>++` - Increases the karma of `<thing>`.
- `<thing>--` - Decreases the karma of `<thing>`.
- `(a thing with spaces)++` - Use parentheses for things with spaces.
- `!karma <thing>` - Checks the current karma score of `<thing>`.

## 3. Quotes

The quote driver stores and retrieves funny or memorable lines.

- `!qadd <quote>` (or `!quote add`, `!add quote`) - Adds a new quote.
- `!quote <regex>` - Finds quotes matching the regular expression.
- `!quote #<ID>` - Retrieves a specific quote by ID.
- `!qdel #<ID>` (or `!quote del`, `!del quote`) - Deletes a quote.

## 4. Reminders and Tells

sp0rkle can remind you of things later or leave messages for absent users.

### Reminders
- `!remind <nick> <message> in <duration>` (e.g., `!remind me coffee in 10m`).
- `!remind <nick> <message> at <time>` (e.g., `!remind me meeting at 2pm`).
- `!remind list` - Lists your pending reminders.
- `!remind del <N>` - Deletes reminder N from your list.
- `!snooze [duration]` - Delays a recently triggered reminder.

### Tells
- `!tell <nick> <message>` (or `!ask`) - Stores a message that the bot will deliver when `<nick>` next speaks or joins.

### Timezones
- `!my timezone is <zone>` - Sets your preferred timezone (e.g., `Europe/London`) for reminders.
- `!forget my timezone` - Clears your timezone setting.

## 5. Decisions

Need help making a choice?
- `!decide <option1> or <option2> or ...` (or `!choose`) - The bot will pick one for you.
- `!rand <range>` - Pick a random number.

## 6. Stats and Seen

- `!seen <nick>` - Check when `<nick>` was last seen.
- `!stats [nick]` (or `!lines`) - Check line count statistics.
- `!topten` (or `!top10`) - Show top 10 most active users.

## 7. Utilities

- `!calc <expression>` - Basic calculator.
- `!date <time/date> [in <zone>]` - Parse and format dates/times.
- `!urbanDictionary <term>` (or `!ud`) - Look up a term on Urban Dictionary.
- `!urlfind <regex>` (or `!urlsearch`) - Search for URLs previously mentioned.
- `!randurl` (or `!random url`) - Displays a random URL from the database.
- `!shorten <url>` - Shortens a URL.
- `!cache <url>` (or `!save`) - Caches a URL.
- `!base <from>to<to> <num>` - Base conversion (e.g., `!base 10to16 255`).
- `!length <string>` - Returns the length of a string.
- `!chr <int>` / `!ord <char>` - Convert between characters and their integer values.
- `!netmask <ip/cidr>` - Calculate netmask information.
- `!markov [nick]` - Generate a random sentence using Markov chains.
- `!markov me` / `!don't markov me` - Opt-in or opt-out of Markov chain learning for your nick.
- `!insult <nick>` - Insult someone at random.

## 8. Integrations

### GitHub
- `!file bug <title>` (or `!report bug`) - Creates a new issue on the bot's GitHub repository.
- `!update bug #<number> <comment>` - Adds a comment to an existing GitHub issue.

### Minecraft
- `!mc set <key> <value>` - Configure Minecraft server integration.

### Push Notifications (Pushbullet)
- `!push enable` / `!push disable` - Toggle push notifications.
- `!push auth <pin>` - Authenticate your Pushbullet account.
- `!push add alias` / `!push del alias` - Manage push notification aliases.

## 9. Bot Management (Authorized Users)

- `!ignore <nick>` / `!unignore <nick>` - Make the bot ignore/unignore a user.
- `!rebuild` - Triggers a bot rebuild and restart.
- `!migrate <state>` - Manually trigger database migration steps.
