## **gigit - A simple CLI tool for fetching and cloning GitHub repositories**

**gigit** is a command-line tool for easily fetching, downloading, and cloning GitHub repositories. It allows users to retrieve repositories by specifying the user, repository name, and commit hash or branch. The tool supports downloading repository archives, cloning repositories, and handling specific subdirectories or tags with simple commands.

#### **Features:**
- **Fetching Repositories:** Supports fetching repositories by user/repo name, including commit hashes, branches, or tags.
- **Subdirectory Fetching:** Ability to fetch specific subdirectories within repositories.
- **Cloning Repositories:** Clone repositories from GitHub with support for shallow clones (`--depth=1`).
- **Versioning Support:** Handles specific commits or version tags (e.g., `v1.0.0`).
- **Error Handling & Retries:** Provides helpful error messages and retries the fetch or clone process if needed.

#### **Commands:**
- `gigit user/repo`: Fetches the latest commit of the specified repository.
- `gigit user/repo/subdir`: Fetches the repository and specific subdirectory.
- `gigit user/repo#commit`: Fetches a specific commit or branch.
- `gigit clone user/repo`: Clones the specified repository.
- `gigit c1 user/repo`: Clones the repository with `--depth=1` for a shallow clone.
- `gigit help`: Displays help information for using the tool.

#### **Installation:**
Install the `gigit` CLI tool using Go:

```sh
go install github.com/nazhard/gigit@latest
```

#### **Example Usage:**
```sh
gigit nazhard/gigit           # Fetch the latest commit of the 'gigit' repository
gigit nazhard/gigit/cmd       # Fetch the 'cmd' subdirectory from 'gigit'
gigit nazhard/gigit#v1.0.0    # Fetch a specific version 'v1.0.0' of 'gigit'
gigit clone nazhard/gigit     # Clone the 'gigit' repository
gigit c1 nazhard/gigit        # Clone the 'gigit' repository with depth=1
```
# Error Handling:
- Invalid commands or arguments are met with helpful error messages, guiding the user on the correct usage and providing suggestions for valid commands 
