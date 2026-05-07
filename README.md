## final-sprint1
### Test plan
- GitHub Actions: job `test` passes on push/PR
- GitHub Actions: job `deploy` runs on tag push and publishes image to Docker Hub
- Docker: `docker pull paulpavlo/final-sprint1:<tag>` and `docker run` starts successfully

CD: deploy runs on tag push (v*).
