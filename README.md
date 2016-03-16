[![Build Status](https://travis-ci.org/GrappigPanda/notorious.svg?branch=devel)](https://travis-ci.org/GrappigPanda/notorious) [![Go Report Card](https://goreportcard.com/badge/github.com/GrappigPanda/notorious)](https://goreportcard.com/report/github.com/GrappigPanda/notorious)
# Notorious
Hello everybody! Notorious aims to be a highy extensible tracker implemented in golang. Right now, Notorious uses Redis to store peer information for quick retrieval and to alleviate some of the burden of storing peer information from the tracker itself. Moreover, some of the core functionalities which Notorious hopes to gain include the following:
```
1. Ratio tracking using a SQL backend (I will be using an ORM layer so that most SQL DBs will be supported).
2. "Complete" Ratioless tracking: no peer information is ever stored
3. "Semi-Ratioless" tracking: user information is stored, but only what is import: grabs/seed life/&c.
4. Automatic Redis docker deployment. I like docker and Redis is my peer-storage of choice and they work well together.
5. Speed and scalability are always in the back of my mind, even if my decisions don't always reflect that
```


There's probably a lot more! Check out my [issues page](https://github.com/GrappigPanda/notorious/issues)


Hey, while you're here why don't you try out NetBSD after testing out my tracker [NetBSD-7.0.arm64.torrent](NetBSD-7.0.arm64.torrent)
(I'm not part of the NetBSD project or anything, I just like the product).



Notorious is a project which I've had a ton of fun learning Go in, but do realize I'm still learning Go so I do make non-idiomatic decisions. If you see anywhere that you think I could improve my code or golang usage, please:
[open an issue](https://github.com/GrappigPanda/notorious/issues/new)
[tweet me](http://twitter.com/GrappigPanda)
[or email me](mailto:ian@ianleeclark.com)
