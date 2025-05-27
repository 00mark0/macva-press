package db

import (
	"log"
	"testing"

	"github.com/00mark0/macva-press/utils"

	"context"

	"github.com/stretchr/testify/require"
)

func createCategoryInteractive(name string) Category {
	category, err := testQueries.CreateCategory(context.Background(), CreateCategoryParams{
		CategoryName: name,
		Slug:         utils.Slugify(name),
	})
	if err != nil {
		log.Fatal(err)
	}
	return category
}

func createCategories(t *testing.T) []Category {
	categoryStrings := utils.RandomCategoryList(10)
	var categories []Category
	for i := 0; i < 10; i++ {
		category := createCategoryInteractive(categoryStrings[i])
		require.NotEmpty(t, category)
		require.Equal(t, categoryStrings[i], category.CategoryName)

		categories = append(categories, category)
	}
	return categories
}

func TestCreateCategory(t *testing.T) {
	categories := createCategories(t)
	require.NotEmpty(t, categories)
	require.Len(t, categories, 10)

	for _, category := range categories {
		require.NotEmpty(t, category)
		require.NotEmpty(t, category.CategoryID)
		require.NotEmpty(t, category.CategoryName)
	}
}

func TestGetCategory(t *testing.T) {
	randomCategoryName := utils.RandomCategory()

	category1, err := testQueries.CreateCategory(context.Background(), CreateCategoryParams{
		CategoryName: randomCategoryName,
		Slug:         utils.Slugify(randomCategoryName),
	})
	require.NoError(t, err)

	category2, err := testQueries.GetCategory(context.Background(), category1.CategoryID)
	require.NoError(t, err)
	require.Equal(t, category1.CategoryID, category2.CategoryID)
	require.Equal(t, category1.CategoryName, category2.CategoryName)
}

func TestGetCategoryByName(t *testing.T) {
	randomCategoryName := utils.RandomCategory()
	category1, err := testQueries.CreateCategory(context.Background(), CreateCategoryParams{
		CategoryName: randomCategoryName,
		Slug:         utils.Slugify(randomCategoryName),
	})
	require.NoError(t, err)

	category2, err := testQueries.GetCategoryByName(context.Background(), category1.CategoryName)
	require.NoError(t, err)
	require.Equal(t, category1.CategoryID, category2.CategoryID)
	require.Equal(t, category1.CategoryName, category2.CategoryName)
}

func TestListCategories(t *testing.T) {
	categories, err := testQueries.ListCategories(context.Background(), 5)
	require.NoError(t, err)
	require.LessOrEqual(t, len(categories), 5)

	for _, category := range categories {
		require.NotEmpty(t, category)
		require.NotEmpty(t, category.CategoryID)
		require.NotEmpty(t, category.CategoryName)
	}
}

func TestUpdateCategory(t *testing.T) {
	randomCategoryName := utils.RandomCategory()

	category1, err := testQueries.CreateCategory(context.Background(), CreateCategoryParams{
		CategoryName: randomCategoryName,
		Slug:         utils.Slugify(randomCategoryName),
	})
	require.NoError(t, err)

	arg := UpdateCategoryParams{
		CategoryID:   category1.CategoryID,
		CategoryName: utils.RandomString(10),
	}

	category2, err := testQueries.UpdateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, category1.CategoryID, category2.CategoryID)
	require.Equal(t, arg.CategoryName, category2.CategoryName)
}

func TestDeleteCategory(t *testing.T) {
	randomCategoryName := utils.RandomCategory()

	category1, err := testQueries.CreateCategory(context.Background(), CreateCategoryParams{
		CategoryName: randomCategoryName,
		Slug:         utils.Slugify(randomCategoryName),
	})

	count1, err := testQueries.ListCategories(context.Background(), 20)
	require.NoError(t, err)
	require.NotZero(t, len(count1))

	deleted, err := testQueries.DeleteCategory(context.Background(), category1.CategoryID)
	require.NoError(t, err)
	require.Equal(t, category1.CategoryID, deleted.CategoryID)

	count2, err := testQueries.ListCategories(context.Background(), 20)
	require.NoError(t, err)
	require.Equal(t, len(count2), len(count1)-1)
}
