# The Enably Roadmap

**NOTE**: I have written this document so that I could share it amongst friends before starting work on Enably to collect feedback. It is no longer being updated, please see [README.md](../README.md) instead, but I'm leaving it here for posterity's sake.

This document roughly describes my vision for Enably. It also describes the layout of the Enably website and justifies the major technical decisions made.

**Note**: This spec doesn't describe the current state of things, but rather the way things are supposed to be in the future. Enably is still in a very early stage of development, and almost none of the features described here are implemented yet. For more details about what's already done and what still needs doing, see the "Execution Plan" section.

TBD means yet to be determined.

## What problems will Enably solve?

Enably lets people find products that fulfill their accessibility needs. Because I'm a blind person, that's the market I'm focusing on at first, but I might expand it to other markets in the future.

In some product categories, even though there are a lot of options to choose from, most choices have major accessibility issues. For example, microwaves, washing machines and induction stoves often have touch controls, which make them difficult to use for blind people. Sometimes these issues can be worked around to some extend, for example by using an associated mobile app, which can have accessibility issues of its own. This problem is not limited to kitchen appliances, it also affects mobile apps and desktop programs, websites and web apps, games, and even music and ham-radio equipment.

Enably aims to be the one-stop-shop for accessibility information. It's going to offer accessibility ratings for products, any necessary workarounds a blind person might need to use them, blindness-specific information such as button layouts, or links to extra software created specifically to enhance accessibility. All content on Enably will be user generated. users will have the ability to add new products or update them when necessary, for example when a new update causes accessibility problems.

## website structure

Products in the Enably directory will be divided into categories. Categories will be hierarchical, i.e. categories might have subcategories, which might have subcategories of their own. A single category may contain either subcategories or products, not both. This way, users don't accidentally add products to top-level categories like "games", adding them to the appropriate subcategories instead. Most categories will have a subcategory called "other" for products which don't fit anywhere else. We internally call the categories that contain products "leaf categories", and categories that contain other categories "branch categories".

Each product will belong to exactly one category, I can't envision a use case where a product should be in more than one, but this might change in the future.

Users will be able to propose new products and suggest updates to existing ones. When suggesting an update, users will need to justify why the update is required. All updates and submissions will need to be reviewed by a moderator, i.e. me for now. Changes to the category structure and category fields will require touching the code, and will be performed via the traditional Github flow.

Each category will have one or more associated fieldsets. A fieldset is a set of fields that describes some aspect of a product. One category can have one or more fieldsets, and one fieldset can be associated with one or more categories. For example, both iOS apps and Android apps will have the basic_app_details fieldset, containing fields like app name, price, description and accessibility rating. However, iOS apps might have a separate fieldset for Voice-Over-specific information, while Android apps might have a fieldset that indicates whether the app is compatible with Talkback, Commentary or both. This design allows us to reuse the same fields across different categories, while still giving us the flexibility to add category-specific fields.

When adding a product, the user will have to fill out all the (non-optional) fields from the fieldsets in its category and all its parent categories. This way, we can add the basic_app_details fieldset to the "apps" category, and that fieldset will apply to all products in all subcategories.

When adding and viewing a product, fields will be divided into sections, one section per fieldset. Some fields may be marked as optional. If a field is empty, it will not be displayed at all when the product is viewed. This can happen if the field is optional and was intentionally left blank, or when the field didn't exist when the product was added.

The following field types will be supported:

- short text, such as product names and short descriptions
- long, multiline text, for product descriptions and detailed notes
- one option from a specified set, for accessibility rating and such
- many options from a specified set, for supported screen readers
- Note, only displayed when editing a product. Provides guidance on what to put in other fields of the fieldset. Not editable by the user in any way.

### User authentication

Users will be identified by an email address and will not have a password. Logging in and signing up will essentially be the same process, using the exact same page. Whenever a login or sign up is required, the user will be asked for their email address, and a link will be emailed to them to verify the attempt. The actual login will occur on the device the link is clicked on, not the device from which the attempt was made. This might be a little inconvenient, but it prevents most phishing attacks. When a login occurs for an account that doesn't exist yet, the user will be asked a few questions and the account will be created. This way of handling user accounts, even though a little unusual, has a few important benefits. It simplifies things on our end, as we don't have to handle two separate processes for logging in and signing up. We also don't have to store any passwords, which minimizes the consequences of any eventual leaks. This way, the only personal data we ever store are the users' email addresses. Password leaks aren't really a concern either, as most email providers have better defenses against them than we ever could. We also free ourselves from having to implement processes for password changes, password resets and a whole "forgot password" flow. This way of authenticating is potentially more confusing to users, but I think that if it is explained well, it should be easy enough to follow. 

The only downside on our end is that we need to get a transactional email provider, but I think we'd have to do that anyway because of email verification and password resets.

When signing up, users will be asked to provide a their name. This will not be a username used for logging in, but rather a display name used on product pages and so on. Users will be encouraged to provide their full names, but a commonly used nickname is acceptable too. They will also be asked for consent to receive occasional emails from us. We will not be sending anything out yet, but adding this field costs us nothing, and it's good to have that option in the future.

## Technical Notes

**NOTE**: This section is fairly technical, non-programmers may skip it entirely.

Enably is divided between the Backend (written in Go) and the frontend (written in React and Javascript). A REST api is used for communication, as I don't think Graph QL would bring us many benefits right now, and it would complicate our stack. The database remains TBD, I'm considering PostgreSQL and SQLite. SQLite is much simpler to maintain and back up, and the expected traffic volume is low enough for SQLite to handle, but I think Postgres is easier to deploy to platforms like render.com (and formerly Heroku) as they use ephemeral containers with no permanent file storage. It really depends on our chosen deployment option though.

The deployment story remains TBD, the options considered are Render, which I haven't played with but people have been saying nice things about, fly.io (same story) and a custom VPS. The advantage of Fly and Render is that I don't have to do dev ops and worry about security updates myself, and they're free at first. The (potential) disadvantage of Render is the inability to use SQLite; I know for a Fact that Fly is big on SQLite and has good support for it. Fly also seems like a small, nice and responsive company, which is important when accessibility feedback is concerned. For now, Fly really seems like the best option, but the decision hasn't been made yet.

To avoid writing too much custom CSS, which requires help from a sighted person, a UI framework (TBD, but probably Bootstrap) will be used. Other options need to be looked at.

Then there's the question of how to actually store product data in the database. I can see two ways, one table per fieldset and one giant table for all fields. The one table per fieldset approach seems a little weird, interpolating table names in code isn't something people do very often, and it might require multiple queries for use cases like getting all fields for a given product. This approach also requires us to write DB migrations whenever fields change. However, if we go the giant table way, we lose all the power of SQL, as the database can no longer enforce data types, check if fields exist or whether fields aren't duplicated etc. It will probably also slow down queries that filter on certain product aspects and make them harder to optimize, if we ever decide to implement that. I'm leaning towards the first approach, but some more reading is still in order. 

The schema containing the category structure, fieldsets and their associated fields will be stored in a file along with the code. The file format is TBD, but I'm considering YAML, TOML and HCL.

The backend architecture will be as follows:

- The `model` package, doesn't depend on anything and everything depens on it. Contains the raw structs that other layers can use, along with some basic methods operating on those structs. No database code.
- The `store` package, used to communicate with the database. Methods correspond to database operations with no extra business logic.
- The `app` package, contains business logic, declares interfaces which `store` implements, no dependencies between `app` and `store` in either direction. This is where most tests should live. How (and whether) to mock the store is TBD.
- The `api` package, declares interfaces for `app` to implement, handles requests, responses and does the required permission checks. 
- The `main` package at `cmd/enably`, which binds all the other stuff together and passes around the right dependencies. Main depends on everything else, but `app`, `store` and `api` don't depend on each other.


## Open Issues

- Write the execution plan section, link to it in the disclaimer at the top.
- Finish the category list.

## Possible future ideas

- length limit for short text
- formatting for multiline fields
- marking fields and options as inactive to make them unusable for new products but display on old products properly.
- an "other" option for single and multi select fields, with a way to put in custom text.
- letting users see unreviewed products and changes if requested.
- letting users see change history
- breadcrumbs for categories on category and product pages.

## Website Layout

### General

Each page will consist of a navigation bar,the main page contents and a footer. The navigation bar will contain the following links:

- Enably, leading to the homepage
- by category
- New Products
- Recently Updated
- Add a product
- Log In (if not logged in already)
- Sign Up (if not logged in already)
- Logged in as \<username\> (not a link)
- log out (if applicable)

The "Log in" and "Sign up" links will actually open the same login/sign up dialog.

The footer will have a copyright and links to the TOS and privacy policy.

### The Home Page

URL: https://enably.me/

The home page will consist of the following elements

- The Enably name
- A short description of what enably is, actual description TBD.
- A "popular categories" section, with a hardcoded selection of categories. The section will have a "see more" link, which will lead directly to the categories page. The actual list of categories displayed TBD.
- A "New Products" section, with about 5 or so elements and a "See more" link leading to the new products page. Listings displayed as product cards, see below.
- A "Recently Updated" section, same as above.

### The Branch Category page.

URLs:

- /by-category/ - for the root category.
- /by-category/category-id-1/category-id-2/category-id-n - for a normal category.

This page is used to display a branch category (including the root category). It should contain the following elements:

- The name of the category in question, or "Browse by category" in the case of the root category.
- The text "this category contains {n} subcategories", except in the case of the root category.
- A list of subcategories.

### The leaf category page

Same URLS as above

This page is used to display the list of products in a given category. Leaf categories, i.e. categories that contain products, have no subcategories.

The page should contain the following elements:

- Category name
- a "sort by" button that opens a menu. Options should include "newest first", "recent updates" and "a to z".
- A list of products in a given category, displayed as product cards. Clicking on a product opens the product details page.

### The product details page

URL: /products/id

Page elements:

- The product name
- A link called "Correct outdated or inaccurate information", leading to the "Edit product" page.
- The full product description
- Extra sections, one per fieldset. Each section contains fields in a given fieldset, one per line. Simple fields use the format "name:  value", for example "Price: free". longer fields, for example accessibility workarounds, are a heading level 3, with information underneath.
- A link saying "Add more details", leading to the "Edit product" page.

### The "Choose a category for this product" page.
URL: "/submit/optional_cateogy_path"

This page will be displayed when clicking the "Submit" link in the navbar.

Page elements:

- The heading "Choose a category for your new product"
- The text "this category contains {n} subcategories", for non-root categories.
- A list of (sub)categories. Clicking a category redirects to this page, but with the category selected.
- a "go back" link, only when in a non-root category.

### The product form page

This page will be displayed when clicking the "Add new product" link on a category page, or when a leaf category is selected on the previous page. It will also be displayed when editing an existing product.

Page elements:

- A heading saying "Add your new product" or "Edit product details"
- The "basic info" section containing fields for name and product description.
- A disclaimers section with the following checkboxes (only when adding a new product).
	- I'm not associated with the creators of this product.
	- I have no financial motivations in adding this product to the Enably directory.
	- I do not expect that adding this product to the Enably directory will bring me any benefits, financial or otherwise.
	- I have not been asked to add this product by a third party.
- A multiline text area saying "Why are you editing this product?" Only when editing products.
- The "Submit" button. In case of form validation errors, the first section with an error gets expanded, the focus is placed there and an alert is announced.

When a submit is successful, we should be returned to the previous page, with the following alert.

"Your product was added/edited. It/ the change  will appear in the Enably directory once it is approved by a moderator."

## Category layout and product fields

This section exists to guide the design of how fields, categories and products are related, as well as to discover possible field types.

### Category: Mobile, Desktop and Web Apps

- name, required, short text
- short summary, what does this app do? Displayed on product pages. Not too short, not too long.
- General app description, not accessibility specific, multiline long text
- usability
	- all parts are usable with all tested screen readers
	- some relatively unimportant parts of the app are inaccessible with one or more screen readers
	- important parts of the app are inaccessible, but the app is still somewhat usable
	- the app is inaccessible to the point of unusability
- ease of use (conditional on the app being accessible at all, so maybe we should merge these two?)
	- the app is easy to use with assistive technologies
	- the app has some usability issues, such as unlabelled buttons or an unintuitive layout
	- the app requires extensive workarounds because of its usability issues
	- the app is completely unusable with assistive technologies
	- not applicable, the app is inaccessible
- Is external software, such as addons or scripts, required to use this app? Screen readers not included.
	- no, no external software exists
	- external software exists to enhance accessibility or provide a better user experience, but it is not required
	- external software greatly helps with making the app accessible
	- external software is absolutely essential to make this app accessible
	- not applicable, the app is inaccessible
- Do the developers of this app care about accessibility?
	- The app is developed with blind users in mind
	- The developers care about accessibility deeply, and that commitment was demonstrated by their actions, not just their words
	- The app has special accessibility features (for example rotor actions) which demonstrate that the developers care somewhat, but issues often slip through
	- the developers seem uninterested
	- No information
- How do the developers respond to accessibility feedback?
	- Responses are clearly written by a human, and reported issues get fixed promptly.
	- Accessibility issues receive standard customer-support responses and reports don't always lead to a fix
	- no information
- extra notes about the accessibility of this app, if any
- useful tricks to enhance accessibility, including external accessibility-enhancing software
- app price:
	- completely free, with or without ads
	- free, but premium features require payment
	- payment is required, but a trial is available
	- no trial is available, a payment is required
- price/currency, if applicable. Provide U.S. Dollars if possible, describe what you pay for and how often, one time purchase / subscription
- related products, versions for other platforms, special clients for accessibility, scripts / NVDA addons, accessibility-enhancing plugins etc.
- platform, Windows / iOS / Android / Web / NVDA Addon / Mac OS / Linux / Apple Watch / TV OS / Play Station / xBox / Roku. If same interface and a11y issues on multiple platforms provide more than one, otherwise make multiple listings. Only include a platform if you've personally tested on it within the last six months.
- Platforms that you haven't tested on, same list as above
- link to download or purchase the product

more categories to be done later.
