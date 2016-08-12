# Configuring Yubikey with Notary/DCT

Notary can be used with Yubikey 4 keys, via a PKCS11 interface when the Yubikey has CCID mode enabled. The Yubikey will be prioritized to store root keys, and will require user touch-input for signing. Note that Yubikey support is included with the Docker Engine 1.11 client for use with Docker Content Trust.

1. Plug in your Yubikey

2. If you haven't already, you should get a PIV library as Yubikey support requires Yubico PIV libraries. 

    The quickest way to do this is to `brew install yubico-piv-tool`
    
    If you choose not to use `brew`, you can get a tarball from [yubico](https://developers.yubico.com/yubico-piv-tool/Releases/), unpackage it and then build it by running `./configure`, `make`, `sudo make install` as described [here](https://developers.yubico.com/yubico-piv-tool/)
    
    
3. Try a pull of an official image with content trust enabled, eg: `DOCKER_CONTENT_TRUST=1 docker -D pull alpine` 

    If your Yubikey is found and configured, you should see the output-

    ```
    ...

    DEBU[0001] Initialized PKCS11 library /usr/local/lib/libykcs11.so and started HSM session
  
    ...
    ```

    If there is no Yubikey detected, you should see the following-

    ```
    ...

    DEBU[0001] Failed to initialize PKCS11 environment: found library /usr/local/lib/libykcs11.so, but initialize error pkcs11: 0x6: CKR_FUNCTION_FAILED

    ...
    ```

    You could also test this out against your notary client by running `notary -s https://notary.docker.io -d ~/.docker/trust -D list docker.io/library/alpine` 
    
    You should see `DEBU[0000] Initialized PKCS11 library /usr/local/lib/libykcs11.dylib and started HSM session` if your Yubikey is configured successfully else if there is an error in your Yubikey setup you will run into `DEBU[0000] No yubikey found, using alternative key storage: loaded library /usr/local/lib/libykcs11.dylib, but no HSM slots found` 
                                                                       

    If your Yubikey was plugged in and your piv tool was configured and you still ran into the issues, you should check your driver configuration and try using `yubico-piv-tool -a status` to troubleshoot further
