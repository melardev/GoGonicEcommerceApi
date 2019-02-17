package seeds

import (
	"github.com/icrowley/fake"
	"github.com/jinzhu/gorm"
	"github.com/melardev/api_shop_gonic/models"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

func randomInt(min, max int) int {

	return rand.Intn(max-min) + min
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func seedAdmin(db *gorm.DB) {
	count := 0
	adminRole := models.Role{Name: "ROLE_ADMIN", Description: "Only for admin"}
	query := db.Model(&models.Role{}).Where("name = ?", "ROLE_ADMIN")
	query.Count(&count)

	if count == 0 {
		db.Create(&adminRole)
	} else {
		query.First(&adminRole)
	}

	adminRoleUsers := 0
	var adminUsers []models.User
	db.Model(&adminRole).Related(&adminUsers, "Users")

	db.Model(&models.User{}).Where("username = ?", "admin").Count(&adminRoleUsers)
	if adminRoleUsers == 0 {

		// query.First(&adminRole) // First would fetch the Role admin because the query status name='ROLE_ADMIN'
		password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		// Approach 1
		user := models.User{FirstName: "AdminFN", LastName: "AdminFN", Email: "admin@golang.com", Username: "admin", Password: string(password)}
		user.Roles = append(user.Roles, adminRole)

		// Do not try to update the adminRole
		db.Set("gorm:association_autoupdate", false).Create(&user)

		// Approach 2
		// user := models.User{FirstName: "AdminFN", LastName: "AdminFN", Email: "admin@golang.com", Username: "admin", Password: "password"}
		// user.Roles = append(user.Roles, adminRole)
		// db.NewRecord(user)
		// db.Set("gorm:association_autoupdate", false).Save(&user)

		if db.Error != nil {
			print(db.Error)
		}
	}
}

func seedUsers(db *gorm.DB) {
	count := 0
	role := models.Role{Name: "ROLE_USER", Description: "Only for standard users"}
	q := db.Model(&models.Role{}).Where("name = ?", "ROLE_USER")
	q.Count(&count)

	if count == 0 {
		db.Create(&role)
	} else {
		q.First(&role)
	}

	var standardUsers []models.User
	db.Model(&role).Related(&standardUsers, "Users")
	usersCount := len(standardUsers)
	usersToSeed := 20
	usersToSeed -= usersCount
	if usersToSeed > 0 {
		for i := 0; i < usersToSeed; i++ {
			password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
			user := models.User{FirstName: fake.FirstName(), LastName: fake.LastName(), Email: fake.EmailAddress(), Username: fake.UserName(),
				Password: string(password)}
			// No need to add the role as we did for seedAdmin, it is added by the BeforeSave hook
			db.Set("gorm:association_autoupdate", false).Create(&user)
		}
	}
}

func seedTags(db *gorm.DB) {
	var tags [3]models.Tag

	db.Where(&models.Tag{Name: "Shoes"}).Attrs(models.Tag{Description: "Shoes for everyone", IsNewRecord: true}).FirstOrCreate(&tags[0])
	db.Where(models.Tag{Name: "Jackets"}).Attrs(models.Tag{Description: "Jackets for everyone", IsNewRecord: true}).FirstOrCreate(&tags[1])
	db.Where(models.Tag{Name: "Jeans"}).Attrs(models.Tag{Description: "Jeans for everyone", IsNewRecord: true}).FirstOrCreate(&tags[2])

	for _, tag := range tags {
		for i := 0; i < randomInt(1, 3); i++ {
			if tag.IsNewRecord {
				db.Create(&models.FileUpload{Filename: randomString(16) + ".png", OriginalName: randomString(16) + ".png",
					FilePath: "/static/images/tags/" + randomString(16) + ".png", FileSize: 2500,
					Tag: tag, TagId: tag.ID})
			}
		}
	}
}

func seedCategories(db *gorm.DB) {
	var categories [3]models.Category
	db.Where(models.Category{Name: "Women"}).Attrs(models.Category{Description: "Clothes for women", IsNewRecord: true}).FirstOrCreate(&categories[0])
	db.Where(models.Category{Name: "Men"}).Attrs(models.Category{Description: "Clothes for men", IsNewRecord: true}).FirstOrCreate(&categories[1])
	db.Where(models.Category{Name: "Kids"}).Attrs(models.Category{Description: "Clothes for kids", IsNewRecord: true}).FirstOrCreate(&categories[2])

	for _, category := range categories {
		for i := 0; i < randomInt(1, 3); i++ {
			if category.IsNewRecord {
				db.Create(&models.FileUpload{Filename: randomString(16) + ".png", OriginalName: randomString(16) + ".png",
					FilePath: "/static/images/categories/" + randomString(16) + ".png", FileSize: 2500,
					Category: category, CategoryId: category.ID})
			}
		}
	}
}

func seedProducts(db *gorm.DB) {
	productsCount := 0
	productsToSeed := 20
	db.Model(&models.Product{}).Count(&productsCount)
	productsToSeed -= productsCount

	if productsToSeed > 0 {
		rand.Seed(time.Now().Unix())
		tags := []models.Tag{}
		categories := []models.Category{}
		db.Find(&tags)
		db.Find(&categories)
		for i := 0; i < productsToSeed; i++ {
			// add a tag and a category for each product
			// faker.RandomInt(0, len(tags))[0]
			tagForProduct := tags[rand.Intn(len(tags))]
			categoryForProduct := categories[rand.Intn(len(categories))]

			product := &models.Product{Name: fake.ProductName(), Description: fake.Paragraph(),
				Stock: randomInt(100, 2000), Price: randomInt(50, 1000),
				Tags: []models.Tag{tagForProduct}, Categories: []models.Category{categoryForProduct}}
			for i := 0; i < randomInt(1, 4); i++ {
				productImage := models.FileUpload{Filename: randomString(16) + ".png", OriginalName: randomString(16) + ".png",
					FilePath: "/static/images/products/" + randomString(16) + ".png", FileSize: uint(randomInt(1000, 23000))}
				product.Images = append(product.Images, productImage)
				db.Set("gorm:association_autoupdate", false).Create(&product)
			}

			/*
				db.Create(&models.FileUpload{Filename: randomString(16) + ".png", OriginalName: randomString(16) + ".png",
					FilePath: "/static/images/tags" + randomString(16) + ".png", FileSize: 2500,
					Tag: tag, TagId: tag.ID})
			*/
		}
	}
}

func seedComments(db *gorm.DB) {
	commentsCount := 0
	commentsToSeed := 20

	allUsers := []models.User{}
	allProducts := []models.Product{}

	db.Model(&models.Comment{}).Count(&commentsCount)
	commentsToSeed -= commentsCount

	if commentsToSeed > 0 {
		rand.Seed(time.Now().Unix())

		db.Find(&allProducts)
		db.Find(&allUsers)

		for i := 0; i < commentsToSeed; i++ {
			userId := allUsers[rand.Intn(len(allUsers))].ID
			productId := allProducts[rand.Intn(len(allProducts))].ID
			sentences := fake.SentencesN(randomInt(2, 6))
			var comment models.Comment

			if rand.Float32() > 0.3 {
				comment = models.Comment{Content: sentences, UserId: userId, ProductId: productId}
			} else {
				// Comment with rating
				comment = models.Comment{Content: sentences, UserId: userId, ProductId: productId, Rating: randomInt(1, 5)}
			}

			db.Set("gorm:association_autoupdate", false).Create(&comment)
		}
	}
}
func seedAddresses(db *gorm.DB) {
	addressesCount := 0
	addressesToSeed := 20

	allUsers := []models.User{}

	db.Model(&models.Address{}).Count(&addressesCount)
	addressesToSeed -= addressesCount

	if addressesToSeed > 0 {
		rand.Seed(time.Now().Unix())
		db.Find(&allUsers)
		var address models.Address

		var city string
		var country string
		var streetAddress string
		var zipCode string
		for i := 0; i < addressesToSeed; i++ {
			city = fake.City()
			country = fake.Country()
			zipCode = fake.Zip()
			streetAddress = fake.StreetAddress()
			address = models.Address{ZipCode: zipCode, StreetAddress: streetAddress, Country: country, City: city}
			if rand.Float32() > 0.4 {
				user := allUsers[rand.Intn(len(allUsers))]
				address.UserId = user.ID
				address.FirstName = user.FirstName
				address.LastName = user.LastName
			} else {
				address.FirstName = fake.FirstName()
				address.LastName = fake.LastName()
			}

			db.Set("gorm:association_autoupdate", false).Create(&address)
		}
	}
}
func seedOrders(db *gorm.DB) {
	ordersCount := 0
	ordersToSeed := 20

	allAddresses := []models.Address{}
	allProducts := []models.Product{}

	db.Model(&models.Order{}).Count(&ordersCount)
	ordersToSeed -= ordersCount

	if ordersToSeed > 0 {
		rand.Seed(time.Now().Unix())
		// Eager load the address's user association
		db.Find(&allAddresses)
		db.Find(&allProducts)

		for i := 0; i < ordersToSeed; i++ {
			address := allAddresses[rand.Intn(len(allAddresses))]

			order := models.Order{TrackingNumber: randomString(16), OrderStatus: randomInt(0, 3), AddressId: address.ID}
			orderItemsForOrder := randomInt(2, 5)
			if rand.Float32() > 0.3 {
				order.UserId = address.UserId
			}
			for j := 0; j < orderItemsForOrder; j++ {
				product := allProducts[rand.Intn(len(allProducts))]
				orderItem := models.OrderItem{ProductName: product.Name, Price: product.Price, Slug: product.Slug,
					ProductId: product.ID,
					UserId:    address.UserId, Quantity: randomInt(1, 8)}

				order.OrderItems = append(order.OrderItems, orderItem)
			}

			db.Set("gorm:association_autoupdate", false).Create(&order)
		}
	}
}

func Seed() {
	db := models.GetDB()
	rand.Seed(time.Now().UnixNano())
	seedAdmin(db)
	seedUsers(db)
	seedTags(db)
	seedCategories(db)
	seedProducts(db)
	seedComments(db)
	seedAddresses(db)
	seedOrders(db)
}
