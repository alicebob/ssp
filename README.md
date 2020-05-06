**Project status**: Please note that this project works, but is otherwise left as-is. I won't maintain this, and I can't provide any support.


# Minimal SSP

A very basic openrtb2.3 SSP used to tests RTB campaigns. It has a few placements
and a few RTB backends configured via a JSON file, and you can see how the
placements behave in a browser. Win-notification URLs are called, but the
system is stateless, so there is no big "money spent" counter.
Basic support for VAST3 video.

# Auction

It's a proper second-price auction, and if there is only a single bid it'll
charge the bid floor for the placement.
