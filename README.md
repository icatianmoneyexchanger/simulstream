# SimulStream Engine

**SimulStream** is a high-performance, containerized RTMP multi-streaming engine designed for game streamers, content creators, and enterprise media broadcasts. It enables zero-overhead simulcasting across YouTube, Twitch, Kick, Facebook Live, TikTok, and custom RTMP/IPv6 endpoints.

---

## 🌟 Key Features

- **Single Ingest, Multi-Egress:** OBS/Broadcaster encodes once to the local SimulStream container; SimulStream handles multi-destination distribution.
- **Rootless & Lightweight:** Built on containerized NGINX RTMP relay engine with native Podman/Docker support.
- **IPv6 Native Ready:** Seamlessly routes across enterprise IPv6 infrastructure and localized mesh networks.
- **Zero Gaming Latency Overhead:** Offloads network streaming overhead from your primary gaming GPU/CPU.
- **Modular Config & Dynamic Push:** Easily add or rotate stream keys and endpoints without restarting the gaming session.

---

## 📁 Architecture Overview

```
 [ OBS Studio / Streamer ]
             │
             │ (Single RTMP Stream -> rtmp://localhost:1935/live/stream)
             ▼
   ┌───────────────────┐
   │ SimulStream Engine│ (Podman / Containerized NGINX-RTMP)
   └─────────┬─────────┘
             ├───────────────────────┼───────────────────────┐
             ▼                       ▼                       ▼
      [ YouTube Live ]        [ Twitch TV ]           [ Kick / Custom ]
```

---

## 🚀 Quick Start

### 1. Build and Run Container
```bash
podman run -d --name simulstream-engine -p 1935:1935 -v ./nginx.conf:/etc/nginx/nginx.conf:ro docker.io/tiangolo/nginx-rtmp
```

### 2. Configure Stream Destinations
Edit `nginx.conf` and update your stream keys under the `live` application section:

```nginx
rtmp {
    server {
        listen 1935;
        chunk_size 4096;

        application live {
            live on;
            record off;

            # YouTube Live
            push rtmp://a.rtmp.youtube.com/live2/YOUR_YOUTUBE_STREAM_KEY;

            # Twitch
            push rtmp://live.twitch.tv/app/YOUR_TWITCH_STREAM_KEY;

            # Kick / Facebook / Custom RTMP
            # push rtmp://.../YOUR_KEY;
        }
    }
}
```

---

## 📄 License
MIT License - Open Source & Commercial Enterprise Ready.
