package product

import (
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"os"
	"post-backend/internal/helper"
	"post-backend/internal/notification"
	"post-backend/internal/setting"
	"strconv"
	"time"
)

func insertProductImages(p ProductServiceImpl, folder string, fileName string, file multipart.File, errChan chan error, productImage *ProductImages, ctx context.Context, tx *sql.Tx) {
	cwd, err := os.Getwd()
	if err != nil {
		errChan <- err
	}

	productImageFileName := fmt.Sprintf("%d-product-%s", time.Now().Unix(), fileName)
	productImagePath := fmt.Sprintf("%s/public/images/%s/%s", cwd, folder, productImageFileName)

	productImage.ImageUrl = productImageFileName
	err = p.ProductRepository.InsertImage(ctx, tx, *productImage)
	if err != nil {
		errChan <- err
		return
	}

	if err := helper.SaveUploadedFile(file, productImagePath); err != nil {
		defer func() {
			if err := os.Remove(productImagePath); err != nil {
				errChan <- err
			}
		}()
		errChan <- err
		return
	}

	errChan <- nil
}

func deleteImage(folder string, imageName string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	productImagePath := fmt.Sprintf("%s/public/images/%s/%s", cwd, folder, imageName)
	fmt.Println("path = ", productImagePath)
	err = os.Remove(productImagePath)
	if err != nil {
		return err
	}
	return nil
}

func CheckProductStock(settingRepository setting.SettingRepository, notificationRepository notification.NotificationRepository, ctx context.Context, tx *sql.Tx, product Product) error {
	minimumStockAlert, err := settingRepository.FindBy(ctx, tx, "minimum_stock_alert")
	if err != nil {
		return err
	}

	minimumStockAlertInt, err := strconv.Atoi(minimumStockAlert.Value)
	if err != nil {
		return err
	}

	if product.Stock < minimumStockAlertInt {
		notification := notification.Notification{
			Title:   "Stock Alert",
			Type:    "stock",
			Message: fmt.Sprintf("The stock for %s is %d, which is below the minimum stock level of %d. Please restock soon.", product.Name, product.Stock, minimumStockAlertInt),
		}
		_, err = notificationRepository.Insert(ctx, tx, notification)
		if err != nil {
			return err
		}
	}

	return nil
}
