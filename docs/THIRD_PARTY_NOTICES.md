# Bundled offline asset and dependency notices

The 40 equipment images and 28 fictional avatars in [the presentation asset manifest](../integration/presentation/assets/MANIFEST.json) are original project-owned illustrations, with editable SVG sources and PNG outputs. [Provenance and reuse permission](../integration/presentation/assets/README.md). They depict fictional identities and stock; no V1 records, external photograph library or real portraits are included.

The unmodified DejaVu Sans font is distributed with its [Bitstream Vera redistribution notice and DejaVu public-domain modification notice](../frontend/src/assets/fonts/DejaVu-LICENSE.txt). The font is bundled for local PDF generation; it is not sold independently or renamed.

The installed, lockfile-pinned PDF libraries preserve upstream copyright and permission text: [jsPDF MIT](licenses/jsPDF-MIT.txt), [AutoTable MIT](licenses/jsPDF-AutoTable-MIT.txt). Existing shared UI dependencies retain [Base UI MIT](licenses/Base-UI-MIT.txt) and [Lucide ISC](licenses/Lucide-ISC.txt). These notices were copied from the installed packages matching the lockfile; original dependency notices remain in installed distributions. npm and Go dependencies remain pinned in their manifests/lockfiles and retain their respective upstream licenses. Dependencies and toolchains are installed/prepared once before offline operation; caches, vendor rehearsals and installed dependency trees are not committed.

Earlier approved application images remain existing assets with their recorded provenance; this checkpoint adds no downloaded photographs. Private owner references and their raw crops are retained locally and excluded from publication.
