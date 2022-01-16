/*
	Bill recommends that every team has a kit project.
	When we talk about project we talk about Single repo of code as well.

	If you are using Git and want a repo in Git the "kit" will have it's own repo.
	The idea of the kit project is that it is our standard library.

	Most of these packages if not all of them are foundational packages.
	So our web support, data base support , logging, config, these types of
	foundational things that all applications that we are using should probably be using.

	See 1.png. Then we also have "Application" project, this is where we can have many
	differnet types of structures.
	Less is more.

	See 1.png is hat they do at Ardan labs.
	There are 4 folders
	1. cmd/
	2. internal/
		2a. platform/
	3. vendor/

	1. cmd/
	This is where our binaries are going to be. We might have a project that might be building multiple
	binaries, multiple services may be we are doing microservices, cli-tools, lambda functions,
	background cron jobs. We don't need one project structure per binary. That would be too much management work.
	We can have easily, 4-5 or may be 8 developers in the project if you organize things properly.
	Multiple binaries will have their own folders inside "cmd".

	2. internal/
	Internal is a rally special folder name because any packages inside of internal will truly be just internal
	to this project and the compiler will protect us.
	If we have two projects and if one is trying to import an internal package from the other.
	Compiler will say "No". We will put business logic, inside of internal.
	So if you got packages that are providing business level logic. Service level logic we will put them
	inside of internal.
	The whole idea right now is that we know what packages are going where,
	based on what the function of these folders are.

	2a. platform/
	The platform package is sitting inside of internal so that we can get the compiler protection as well.
	No body can import those outside of project.
	The platform folders are like the kit folders, they are not the kit folders.
	These are the foundational packages. WE canname thse folders whatever we want but at Ardan labs they use these
	names.
	Some platform folders can be cfg/, log/, etc.

	3. vendor/
	vendor/ is part of the language tooling, it is third part.
	Go has modules, and modules support both vendoring and nonvendoring.
	If you are not going to be using vendoring with modules then you will be relying on a couple of
	files that you will have to keep in this project. Then you will have to rely on social contract that
	are with those third party packages. We don't want ot rely on social contracts
	to make sure that we have a reproducible build.

	Bill use go modules for populating the vendor folders, he also uses go modules to vendor all the source code
	and keep everything in one spot unless it's not reasonable to do so.
	Project like Kubernetes, we cannot vendor on a project like kubernetes.

	(I didn't understand this)
	The kit packages will become vendored packages in our application.
	We will lock those versions in, in not just on the social contract. B
	ut also byy locking al the code in one place.

	------------------------------------------------------------------------------------------

	The whole idea of project structure is to have a mental model around things.

	We do "go mod" to vendor the third party packages.
	we will do "go mod init", "go mod tidy" and then we will do "go mod vendor"

	-----------------------------------------------------------------------------------------
	/cmd is where the binaries are, we will lay a folder in under command with the name of the binary with
	"main.go" under that.

	/internal is where all the business and service logic is going to be.
	---------------------------------------------------------------------------------------------

	Say if we want to build a package that needs to go talk to a server then we will build it under
	platform since it is foundational.

	Bill's suggestion - Inside of internal and platform we should lay these packages flat.
	Unless there's a real relationship on import dependencies. We shouldn't put
	folders inside folders.
	So we should definitely leave the folders inside internal and inside of platform flat, unless there's a
	clear relationship between folders and we will use hierarchy for that.

	These are some dependency choice validations that we can do,
	https://github.com/ardanlabs/gotraining/blob/f4355fce6fb0a161c7d01e39f166065085a26b6a/topics/go/design/packaging/README.md#validation

	One thing good with Ardan labs project structure is that the naming convention is good,
	we are going to say is that the imports can happen down the project tree and not up.

	That is packages inside cmd can import anything. Packages inside of internal can only import
	internal and platform.
	Packages inside of platform can only import from platform and also say, vendor.

	Then in code review if we see a platform package importing something directly inside of internal, the code review
	will stop. The import is going in the wrong direction, that import is going to hurt you later on and not help.

	So we can always be validating these dependency choices and where things are.
	-----------------------------------------------------------------------------------------------------

	Policy is defining dependency choices we make, what databases are we using. What logger we using.
	What configuration are we using.

	Some of the policies we can say are:

	Validate the policies being imposed.

	1. Kit, internal/platform/
		1a. Packages inside "Kit, internal/platform/" should NOT allowed to set policy about any application concerns.
			As these folder locations are for packages that have to have highest level of portability. Highest levels of
			reusability. So if we have a package inside of "kit" that is importing a logger then
			if we don't want to use the logger we can't use that package either. BEcause that package is set for
			policy on logging or policy on config. Packages inside "kit and platform", any sought of import outside of
			the standar library should be double checked.
		1b. NOT allowed to log, but access to trace information must be decoupled.
		1c. Configuration and runtime changes must be decoupled.
	2. Retrieving metric and telemetry values must be decoupled.
		2a. cmd/, internal/
		2b. Allowed to set policy about any application concerns.
		2c. Allowed to log and handle configuration natively.

	Packages in /cmd are at the application level, /internal is the business level.
	These have a lot more flexibility on what they can import as compared to the other packages.





*/