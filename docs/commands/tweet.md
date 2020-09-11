# Tweet

> Listening to users tweet and repost it in the channel

**Trigger**: `s.tweet`

## Stream

> Track a specific user's tweet and post it in the channel whenever there is a new tweet

### Format

`s.tweet stream <twitter-username>`

### Example

`s.tweet stream @username`

## Stop

> Stop tracking a specific user's tweet

### Format

`s.tweet stop-stream <twitter-username>`

### Example

`s.tweet stop-stream @username`

## List

> Get a list of tracked user in the channel

`s.tweet list-stream`