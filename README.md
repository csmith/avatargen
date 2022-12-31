# Stable Diffusion avatar generator and server

Contains two services: `generate` that monitors directories and if they are
low on images, spins up a paperspace instance, feeds a bunch of prompts to
a REST API there to generate images, and saves them to disk. `serve` handles
serving of the avatars with various constraints to ensure some users see
"unique" for-their-eyes-only avatars.

For more information see this blog post: https://chameth.com/infinite-avatars/
