# Image Inbox

Drop raw screenshots and images here for processing.

## How it works

1. Place images in this folder
2. Invoke `@process-planton-ai-images` in Cursor with context about what the images are
3. The tool compresses, renames, organizes, and uploads them to R2
4. Processed images are removed from this folder

## What to include in context

When invoking the rule, describe:
- What feature or page the images relate to
- What the images show
- Any specific naming preferences

## Example

```
@process-planton-ai-images

I've added 3 screenshots showing the deployment pipeline UI in Service Hub.
These are for the deployment-environments documentation.
```

## Notes

- This folder is gitignored (except this README and .gitignore)
- Raw images won't be committed to the repository
- Images are processed, compressed, and uploaded to R2 before being placed in `content/assets/images/`

