# Reminders and Tells

sp0rkle can remind you of things later or leave messages for absent users.

## Reminders
- `!remind <nick> <message> in <duration>` (e.g., `!remind me coffee in 10m`).
- `!remind <nick> <message> at <time>` (e.g., `!remind me meeting at 2pm`).
- `!remind list` - Lists reminders set by or for your nick.
- `!remind del <N>` - Deletes (previously listed) reminder N.
- `!snooze [duration]` - Resets the previously-triggered reminder.

## Tells
- `!tell <nick> <msg>` - Stores a message for the (absent) nick.
- `!ask <nick> <msg>` - Same as `!tell`.

## Timezones
- `!my timezone is <zone>` - Sets a local timezone for your nick.
- `!forget my timezone` - Unsets your local timezone.
