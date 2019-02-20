# Introduction
This is one of my E-commerce API app implementations. It is written in Golang using go-gonic web framework..
This is not a finished project by any means, but it has a valid enough shape to be git cloned and studied if you are interested in this topic.
If you are interested in this project take a look at my other server API implementations I have made with:

- [Node Js + Sequelize](https://github.com/melardev/ApiEcomSequelizeExpress)
- [Node Js + Bookshelf](https://github.com/melardev/ApiEcomBookshelfExpress)
- [Node Js + Mongoose](https://github.com/melardev/ApiEcomMongooseExpress)
- [Django](https://github.com/melardev/DjangoRestShopApy)
- [Flask](https://github.com/melardev/FlaskApiEcommerce)
- [Java EE Spring Boot and Hibernate](https://github.com/melardev/SBootApiEcomMVCHibernate)
- [Ruby on Rails](https://github.com/melardev/RailsApiEcommerce)
- [AspNet Core](https://github.com/melardev/ApiAspCoreEcommerce)
- [Laravel](https://github.com/melardev/ApiEcommerceLaravel)

The next projects to come will be:
- Elixir with phoenix and Ecto
- AspNet MVC 6
- Java EE with Jax RS with jersey
- Java EE with Apache Struts 2
- Spring Boot with Kotlin
- Go with Gorilla and Gorm
- Go with Beego
- Laravel with Fractal and Api Resources
- Flask with other Rest Api frameworks such as apisec, flask restful

## WARNING
I have mass of projects to deal with so I make some copy/paste around, if something I say is missing or is wrong, then I apologize
and you may let me know opening an issue.

# Getting started
1. go get https://github.com/melardev/ApiEcomGoGonic
1. Change the .env.example as you need(see warning below)
1. Rename .env.example to .env
1. Seed the database passing "create seed" as arguments to the app(read main.go to understand what I mean)

## WARNING
The recommended database to use is Postgresql, the other database backends may not work as expected.
Unfortunately the MySQL does not work as expected, for example the BeforeSave Hook for User is not able to retrieve
the Role model if using MySQL, the same code does work if SQLite, it is weird, because the SQL query generated is valid and it
returns a row, but somehow the driver is not able to map it to the user.

# Features
- Authentication / Authorization
- JWT middleware for authentication
- Multi file upload
- Database seed
- Paging with Limit and Offset using GORM (Golang ORM framework)
- CRUD operations on products, comments, tags, categories, orders
![Fetching products page](./github_images/postman.png)
- Orders, guest users may place an order
![Database diagram](./github_images/db_structure.png)

# What you will learn
- Golang
- Golang Go-Gonic web framework
- JWT
- Controllers
- Middlewares
- JWT Authentication
- Role based authorization
- GORM
    - associations: ManyToMany, OneToMany, ManyToOne
    - virtual fields
    - Select specific columns
    - Eager loading
    - Count related association
    
- seed data
- misc
    - project structure

# Understanding the project
The project is meant to be educational, to learn something beyond the hello world thing we find in a lot, lot of 
tutorials and blog posts. Since its main goal is educational, I try to make as much use as features of APIs, in other
words, I used different code to do the same thing over and over, there is some repeated code but I tried to be as unique
as possible so you can learn different ways of achieving the same goal.

Project structure:
- models: Mvc, it is our domain data.
- dtos: it contains our serializers, they will create the response to be sent as json. They also take care of validating the input(feature incomplete)
- controllers: well this is the mvC, they receive the request from the user, they ask the services to perform an action for them on the database.
- seeds: contains the file that seeds the database.
- static: a folder that will be generated when you create a product or tag or category with images
- services: contains some business logic for each model, and for authorization
- middlewares: it contains middlewares(golang functions) that are triggered before the controller action, for example, a middleware which
reads the request looking for the Jwt token and trying to authenticate the user before forwarding the request to the corresponding controller
action

# TODO
- Add model constraints such as not null
- Refactor the seeding with http://gorm.io/docs/query.html#Select
- Global Application Error handling
- Can't Preload field errors:
    - Get comment details http://127.0.0.1:8080/api/products/:slug/comments/:id triggered in services.FetchCommentById
    - Get My Orders http://localhost:8080/api/orders triggered with services.FetchOrdersPage
- Security, validations, file upload
- Delete FileUpload if associated tag, category or product deleted
- Delete Files if tag, category, product fail to be saved
- Use pointers as function parameters instead of passing them by value as I did in many

# Resources
- [Go-Gonic](https://github.com/gin-gonic/gin) Awesome golang based web framework
- [GORM]()
- [CORS gin's middleware](https://github.com/gin-contrib/cors)