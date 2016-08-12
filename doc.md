# Common content trust actions

##Admin actions

####Situation: 

Rotate the root/targets key (and snapshot if kept locally)

> notary key rotate example.com/an_image targets


####Situation: Admin thinks that the passphrase for one of his keys has been compromised

Change the passphrase for a local key

> notary key passwd <keyID>

####Situation:

Rotate the remote snapshot/timestamp key

> notary key rotate example.com/an_image snapshot -r

####Situation: A contributor leaves project and the admin wants their key should be purged from any roles it could sign 

Remove a key from all delegation roles

> notary delegation purge example.com/an_image <keyID>

####Situation:

Rotate a key in a delegation role

> notary delegation add example.com/an_image targets/qa_team <newCert.pem>
> notary delegation remove example.com/an_image targets/qa_team <oldKeyID>

####Situation: For whatever reason, want to remove everything related to a certain Notary repo

Delete all trust data for a repo

> notary delete example.com/an_image --remote # omit the remote flag if you want data to persist on server

##User actions

####Situation: You choose to retract your trust in an image

Remove your signature for an image

> notary remove -p example.com/an_image v1.0 --r targets/qa_team

####Situation: Add your signature to your role's signed images and mark a role to be re-signed on next publish

Add your signature to your role's signed images

> notary witness -p example.com/an_image targets/qa_team

#Troubleshooting Notary issues
- View debug logs

    Notary provides a `-D` flag which prints debug logs for typical commands. These logs help you understand the issues that Notary is running into enabling you to pinpoint them better
    
    For example: `notary -s https://notary.docker.io -d ~/.docker/trust list docker.io/library/alpine` may give you a `* fatal: client is offline` error if your default config is set to trust the default root ca from which you may not be able to tell much. Now if you run the command with a `-D` flag, you will get a large amount of debug logs which contains the line
    
    `ERRO[0000] could not reach https://notary.docker.io: Get https://notary.docker.io/v2/: x509: certificate signed by unknown authority` which may help you pinpoint the issue
    
- Un-staging/resetting staged changes

    You can view staged changes by running `./bin/notary status examplegun`
    
    Should you want to un-stage (reset) these changes, you can use `./bin/notary status examplegun --reset`

- Running into `x509: certificate signed by unknown authority` 

    This is likely the case because you have in your Notary configurations, specified a root certificate to specify a `root certificate authority` which doesn't match with the signature from your server

    If you are using the Notary server `https://notary.docker.io`, you can simply specify a blank config file, config.json > an empty json "{}" and run something like `notary -s -c config.json https://notary.docker.io -d ~/.docker/trust list docker.io/library/alpine` where config.json is an empty json

- Other common issues

    Space for other issues that pop up 
    
    Ideas: 
    - docker-compose up issues so getting an offline client issue

- Fixing your Yubikey setup

    This unit is for users that have plugged in their Yubikey and [configured it](http://linktoconfig.doc), have configured `yubico-piv-tool` and are yet running into `DEBU[0000] No yubikey found, using alternative key storage: loaded library /usr/local/lib/libykcs11.dylib, but no HSM slots found` 
    
    