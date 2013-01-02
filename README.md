A new website, meant to supersede <http://kenjitakahashi.github.com>.

In case you didn't notice, it's hosted @ <http://kenji.sx>.

changelog
---------
0.1.0:
- first public release

0.1.1:
- removed res.redirect hack-in

0.2.0:
- updated jade to 0.27.x
- removed WIP from title
- refactored template(s)
- added notice block

0.2.1:
- fixed small bug in left template (unnoticable)
- fixed large bug in db wrapper (crash-causer)

0.3.0:
- changed notice style
- updated stylus to 0.29.x
- updated highlight.js to 7.1.x
- updated mongoose to 3.0.x
- dropped date-ext dependency
- re-done calendar
- removed unnecessary images

0.4.0:
- updated mongoose to 3.3.x
- updated highlight.js to 7.3.x
- updated stylus to 0.30.x
- updated coffee-script to 1.4.x
- moved to express (tah-dah)
- updated connect to 2.6.x
- embedding canvas
- raw display
- added Google Analytics
- fixed: posts showdown should reset to top when changing date/tag
- removed connect-assets in favour of custom-baked static assets handler

0.4.1:
- moved from Google Analytics to Mixpanel
- updated mongoose to 3.4.x
- updated stylus to 0.31.x
- updated connect to 2.7.x

0.4.2:
- storing assets also in db

0.4.3:
- updated mongoose to 3.5.x
- updated clean-css to 0.9.x
- added mixpanel badge (25000 is not enough it seems)

0.5.0:
- more descriptive mixpanel posts tracking
- RSS feed
- fixed some minor bugs
- serving in-db assets dynamically if not cached on startup
- updated right side a bit

todo
----
* comments
* stylesheet refactor
* scroll lists with wheel
* fix: horizontal code scroll shouldn't affect line numbers
* showing images
* websocketize
* hiding side-menus
