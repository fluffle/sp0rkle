# Factoid Driver

Factoids allow the bot to remember and recall information.

## Adding Factoids
There are two ways to teach the bot something (must address the bot by nick):
- `botNick: <key> := <value>` - Stores `<value>` for `<key>`.
- `botNick: <key> :is <value>` - Stores `<key> is <value>` for `<key>`.

## Retrieving Factoids
- Simply type the `<key>` in the channel, and the bot will respond with a random value associated with that key.
- `botNick: literal <key>` - Prints all known values for a key. If there are many, the bot will ask you to do this in a private message.

## Managing Factoids
These commands operate on the **last factoid triggered** in the channel:
- `botNick: forget that` or `delete that` - Deletes the last triggered value.
- `botNick: replace that with <new value>` - Replaces the last value.
- `botNick: that =~ /regex/replacement/` - Edit the last factoid value using regex.
- `botNick: chance of that is <N>%` - Sets the probability (0-100%) that the bot will respond to this key.
- `botNick: fact info` - Displays metadata about the last factoid (creator, time, etc.).
- `botNick: fact search <regex>` - Search for factoid keys matching a regex.

## Factoid Variables
You can use variables in factoid values:
- `$nick`: The nick of the person who triggered the factoid.
- `$chan`: The channel it was triggered in.
- `$date` / `$time`: Current date or time.
- `$user` / `$host`: Ident or host of the triggerer.
