# envtpl
_a port of https://github.com/andreasjansson/envtpl to Go_

```
$ envtpl something.conf.tpl
```

Add to your `Dockerfile`:

```dockerfile
ENV ENVTPL_VERSION 0.1.0
RUN set -x \
    && apt-get update && apt-get install -y --no-install-recommends ca-certificates wget && rm -rf /var/lib/apt/lists/* \
    && dpkgArch="$(dpkg --print-architecture | awk -F- '{ print $NF }')" \
    && wget -O /usr/local/bin/envtpl "https://github.com/mattrobenolt/envtpl/releases/download/$ENVTPL_VERSION/envtpl-linux-$dpkgArch" \
    && wget -O /usr/local/bin/envtpl.asc "https://github.com/mattrobenolt/envtpl/releases/download/$ENVTPL_VERSION/envtpl-linux-$dpkgArch.asc" \
    && export GNUPGHOME="$(mktemp -d)" \
    && gpg --keyserver ha.pool.sks-keyservers.net --recv-keys D8749766A66DD714236A932C3B2D400CE5BBCA60 \
    && gpg --batch --verify /usr/local/bin/envtpl.asc /usr/local/bin/envtpl \
    && rm -r "$GNUPGHOME" /usr/local/bin/envtpl.asc \
    && chmod +x /usr/local/bin/envtpl \
    && envtpl --help \
    && apt-get purge -y --auto-remove ca-certificates wget
```
