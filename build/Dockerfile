FROM golang:1.24-bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
      libsdl2-dev \
      libsdl2-ttf-dev \
      libsdl2-image-dev \
      libsdl2-gfx-dev \
      libgtk-3-dev \
      libgtk-3-bin \
      dpkg-dev           \
      libgdk-pixbuf2.0-bin  \
       libglib2.0-bin        \
      librsvg2-common  \
      patchelf \
      librsvg2-dev \
      libgdk-pixbuf-2.0-dev \
      desktop-file-utils \
      file \
      ca-certificates \
      wget \
    && rm -rf /var/lib/apt/lists/*

RUN wget -qO /usr/local/bin/linuxdeploy \
      https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous/linuxdeploy-x86_64.AppImage \
 && wget -qO /usr/local/bin/linuxdeploy-plugin-gtk.sh \
      https://raw.githubusercontent.com/linuxdeploy/linuxdeploy-plugin-gtk/master/linuxdeploy-plugin-gtk.sh \
 && chmod +x /usr/local/bin/linuxdeploy /usr/local/bin/linuxdeploy-plugin-gtk.sh

ENV APPIMAGE_EXTRACT_AND_RUN=1
 
WORKDIR /src