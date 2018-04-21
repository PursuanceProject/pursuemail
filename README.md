# PursueMail

## Dev

Used https://mailcatcher.me/ for SMTP testing

Important note: There are a bunch of hardcoded values at the moment.


## Example API Calls

### Map Email Address to (Random) UUID

```
curl -i localhost:9080/api/v1/email -d {"email": "spam@pursuanceproject.org"}
```


### Send Emails

In the below examples, the emails sent to users will be encrypted if
and only if their PGP keys are found in `~/.gnupg/pubring.gpg`,
otherwise they will be sent unencrypted.

If you want to tell PursueMail to only send an email if it is sent in
encrypted form, add `"secure_only": true` to the top level of the JSON
POST body when doing any of the following API calls.


#### Send Email to One User by (UU)ID

```
curl -i localhost:9080/api/v1/email/ec348de2-2430-46d6-9ed7-f65b12a4a75a/send -d '{"email_data": {"from": "team@pursuanceproject.org", "subject": "4 tasks due today!", "body": "4 tasks due today: ..."}}'
```

#### Send Bulk Email (to Multiple Users) by their Email Addresses

```
curl -i localhost:9080/api/v1/email/bulksend -d '{"emails": ["activist1@riseup.net", "activist2@riseup.net"], "email_data": {"from": "team@pursuanceproject.org", "subject": "2 tasks due today!", "body": "2 tasks due today in pursuance #827: ..."}}'
```

#### Send Bulk Email (to Multiple Users) by their (UU)IDs

```
curl -i localhost:9080/api/v1/email/bulksend -d '{"ids": ["ec348de2-2430-46d6-9ed7-f65b12a4a75a", "451724a2-ddb8-4fd9-8336-819316c6019a"], "email_data": {"from": "team@pursuanceproject.org", "subject": "3 tasks due today!", "body": "3 tasks due today in pursuance #827: ..."}}'
```

#### Send _Definitely-encrypted_ Email

Same as these above examples, but add `"secure_only": true` at the top
level.


## TODOs

- [ ] Create a go client library
- [ ] Audit error messages, make sure nothing sensitive is being revealed
- [ ] Better handling of HTML vs. Text emails
- [ ] Support an "Email Settings" page where users can unsubscribe.
- [ ] Better bounce support. (If we spam a non existant email, we are likely to get marked as a spambot).
- [ ] Support a DELETE option? a PUT option?


## Potential Problem Areas

* Currently the sending endpoints hold onto the connection until all
emails are sent.  May want to move this to a background worker with
callback support.
