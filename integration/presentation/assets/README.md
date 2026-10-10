# Offline presentation illustration assets

These68 assets are original project-owned vector illustrations:40 food-service catalog images and28 fictional profile avatars. They depict invented demonstration stock and illustrated people; they are not institutional inventory photographs or real student/faculty portraits. No external photograph library, network download or V1 record was used.

`generate-assets.py` writes the SVG source artwork and equipment metadata. `frontend/scripts/render-presentation-assets.mjs` uses a local Chromium canvas to rasterize those original vectors to256×256 PNGs. `MANIFEST.json` records source/output SHA-256, sizes and dimensions. Runtime storage validates the PNGs through the ordinary catalog/profile image service and keeps them protected in the isolated PostgreSQL database; no remote image URL is required.

The repository owner may use, adapt and redistribute this original artwork with the project. Catalog imagery for live operational equipment remains operator-managed under the existing equipment image rules. These illustrations do not certify real physical assets. They are independent of returns: eLabTrack does not upload/store/accept photographic return evidence.
