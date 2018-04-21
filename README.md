# pursuemail

### Dev

Used https://mailcatcher.me/ for SMTP testing

Important note: There are a bunch of hardcoded values at the moment.

### TODOs

- [ ] Create a go client library
- [ ] Audit error messages, make sure nothing sensitive is being revealed
- [ ] Better handling of HTML vs. Text emails
- [ ] Support an "Email Settings" page where users can unsubscribe.
- [ ] Better bounce support. (If we spam a non existant email, we are likely to get marked as a spambot).
- [ ] Support a DELETE option? a PUT option?

### Potential Problem Areas

* Currently the sending endpoints hold onto the connection until all
emails are sent.  May want to move this to a background worker with
callback support.
