package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
	"szi-registry/csvimport"
	"szi-registry/export"
	"szi-registry/models"
	"szi-registry/services"
	"szi-registry/utils/ui"
	"time"
)

// отображает окно с реестром СЗИ
func ShowSZIRegistryWindow(myApp fyne.App, username string, db *gorm.DB) {
	myWindow := myApp.NewWindow("Реестр СЗИ от НСД — Реестр")
	myWindow.Resize(fyne.NewSize(1400, 600))

	// Получаем ID пользователя по имени
	user, err := services.GetUserByUsername(db, username)
	if err != nil || user == nil {
		// Если не удалось получить пользователя, используем 0
		user = &models.User{ID: 0}
	}

	// Параметры пагинации
	pageSize := 10
	currentPage := 0
	var allRecords []models.SZIRecord


	// Загрузка данных из базы данных для конкретного пользователя
	records, err := services.GetUserRecords(db, uint(user.ID))
	if err != nil {
		records = []models.SZIRecord{}
	}

	// Создаем таблицу
	var table *widget.Table
	table = widget.NewTable(
		func() (int, int) {
			return 0, 12 
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			return label
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
		})

	pageInfoLabel := widget.NewLabel("")

	// Устанавливаем начальные ширины столбцов
	initialColumnWidths := []float32{250, 120, 140, 120, 120, 100, 150, 120, 120, 120, 120, 100}
	for i, width := range initialColumnWidths {
		if i < len(initialColumnWidths) {
			table.SetColumnWidth(i, width)
		}
	}

	// для обновления информации о странице
	updatePageInfo := func() {
		totalRecords := len(allRecords)
		if totalRecords == 0 {
			pageInfoLabel.SetText("Страница 1 из 1")
		} else {
			totalPages := (totalRecords + pageSize - 1) / pageSize
			pageInfoLabel.SetText(fmt.Sprintf("Страница %d из %d", currentPage+1, totalPages))
		}
	}

	// обновления данных в таблице с учетом страниц
	updateTableData := func(data []models.SZIRecord) {
		// Обновляем все записи
		allRecords = data

		totalPages := 1
		if len(data) > 0 {
			totalPages = (len(data) + pageSize - 1) / pageSize
		}
		if currentPage >= totalPages && totalPages > 0 {
			currentPage = totalPages - 1
		} else if totalPages == 0 {
			currentPage = 0
		}

		startIndex := currentPage * pageSize
		endIndex := startIndex + pageSize
		if endIndex > len(data) {
			endIndex = len(data)
		}

		// Ограничиваем данные для текущей страницы
		pageRecords := data[startIndex:endIndex]

		// Обновляем таблицу
		table.Length = func() (int, int) {
			return len(pageRecords), 12 // строки: записи из БД, 12 столбцов
		}
		table.UpdateCell = func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			if id.Row < len(pageRecords) {
				if id.Col < 11 {
					// Обычные данные
					record := pageRecords[id.Row]
					var text string
					switch id.Col {
					case 0:
						text = record.Name
					case 1:
						text = record.Type
					case 2:
						text = record.CertNumber
					case 3:
						text = record.CertIssueDate.Format("2006-01-02")
					case 4:
						text = record.CertExpiryDate.Format("2006-01-02")
					case 5:
						text = ui.GetRecordStatus(record)
					case 6:
						text = record.Manufacturer
					case 7:
						text = record.SoftwareVersion
					case 8:
						text = record.Purpose
					case 9:
						text = record.DeploymentType
					case 10:
						text = record.ClassProtection
					}
					label.SetText(text)
				} else if id.Col == 11 {
					// Столбец действий
					label.SetText("...")
					label.Alignment = fyne.TextAlignCenter
				}
			} else {
				label.SetText("")
			}
		}
		table.Refresh()

		// Обновляем информацию о странице
		updatePageInfo()
	}


	// Обработка нажатий на таблицу
	table.OnSelected = func(id widget.TableCellID) {
		if id.Col == 11 { // Столбец действий
			// Вычисляем индекс записи с учетом пагинации
			startIndex := currentPage * pageSize
			recordIndex := startIndex + id.Row

			if recordIndex >= 0 && recordIndex < len(allRecords) {
				selectedRecord := allRecords[recordIndex]

				// Показываем диалог выбора действия
				actionMenu := widget.NewPopUpMenu(
					&fyne.Menu{
						Label: "Действия",
						Items: []*fyne.MenuItem{
							{Label: "Редактировать", Action: func() {
								ShowEditSziWindow(myApp, selectedRecord, db)
							}},
							{Label: "Удалить", Action: func() {
								// Диалог подтверждения удаления
								var deleteDialog *widget.PopUp

								deleteDialog = widget.NewModalPopUp(
									container.NewVBox(
										widget.NewLabel("Вы уверены, что хотите удалить запись '"+selectedRecord.Name+"'?"),
										container.NewHBox(
											widget.NewButton("Да", func() {
												err := services.DeleteSziRecord(db, uint(user.ID), selectedRecord.ID)
												if err != nil {
													// Показать ошибку
													errorDialog := widget.NewModalPopUp(
														widget.NewLabel("Ошибка при удалении записи: "+err.Error()),
														myWindow.Canvas(),
													)
													errorDialog.Show()
												} else {
													// Обновить окно
													ShowSZIRegistryWindow(myApp, username, db)
													myWindow.Close()
												}
												deleteDialog.Hide()
											}),
											widget.NewButton("Нет", func() {
												deleteDialog.Hide()
											}),
										),
									),
									myWindow.Canvas(),
								)
								deleteDialog.Show()
							}},
						},
					}, myWindow.Canvas())

				// Показываем меню в позиции курсора
				actionMenu.Show()
			}
		}
	}

	// Заголовки столбцов
	headers := []string{"Наименование СЗИ", "Тип СЗИ", "№ сертификата", "Дата выдачи", "Срок действия", "Статус", "Производитель", "Версия ПО", "Назначение", "Тип развертывания", "Класс защищенности", "Действия"}
	columnCount := len(headers)

	// Создаём таблицу-заголовок с такой же структурой, как основная таблица
	headerTable := widget.NewTable(
		func() (int, int) {
			return 1, columnCount // Одна строка, количество столбцов равно количеству заголовков
		},
		func() fyne.CanvasObject {
			label := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
			label.Alignment = fyne.TextAlignCenter
			return label
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label, ok := cell.(*widget.Label)
			if !ok || id.Col >= len(headers) {
				return
			}
			label.SetText(headers[id.Col])
		})

	// Применяем те же ширины столбцов, что и к основной таблице
	for i, width := range initialColumnWidths {
		if i < columnCount {
			headerTable.SetColumnWidth(i, width)
		}
	}

	// Делаем заголовок не кликабельным
	headerTable.OnSelected = func(id widget.TableCellID) {
		headerTable.UnselectAll()
	}

	// Создаем поля для фильтрации
	nameFilter := widget.NewEntry()
	nameFilter.SetPlaceHolder("Поиск по наименованию...")
	typeFilter := widget.NewSelectEntry([]string{"СКЗИ", "ОС", "СЗИ КС", "СЗИ СКЗИ", "Другое"})
	typeFilter.SetPlaceHolder("Тип СЗИ...")
	statusFilter := widget.NewSelectEntry([]string{"Актуально", "Просрочено", "Снято"})
	statusFilter.SetPlaceHolder("Статус...")
	locationFilter := widget.NewEntry()
	locationFilter.SetPlaceHolder("Место установки...")
	manufacturerFilter := widget.NewEntry()
	manufacturerFilter.SetPlaceHolder("Производитель...")
	purposeFilter := widget.NewSelectEntry([]string{"АС", "ИВК", "Универсальное"})
	purposeFilter.SetPlaceHolder("Назначение...")
	deploymentTypeFilter := widget.NewSelectEntry([]string{"Клиент-сервер", "Автономное", "АПК", "Виртуальное"})
	deploymentTypeFilter.SetPlaceHolder("Тип развертывания...")
	classProtectionFilter := widget.NewSelectEntry([]string{"1", "2", "3А", "4", "5"})
	classProtectionFilter.SetPlaceHolder("Класс защищенности...")

	// Функция для применения фильтров
	applyFilters := func() {
		filters := make(map[string]interface{})

		if nameFilter.Text != "" {
			filters["name"] = nameFilter.Text
		}
		if typeFilter.Text != "" {
			filters["type"] = typeFilter.Text
		}
		if statusFilter.Text != "" {
			filters["status"] = statusFilter.Text
		}
		if locationFilter.Text != "" {
			filters["location"] = locationFilter.Text
		}
		if manufacturerFilter.Text != "" {
			filters["manufacturer"] = manufacturerFilter.Text
		}
		if purposeFilter.Text != "" {
			filters["purpose"] = purposeFilter.Text
		}
		if deploymentTypeFilter.Text != "" {
			filters["deployment_type"] = deploymentTypeFilter.Text
		}
		if classProtectionFilter.Text != "" {
			filters["class_protection"] = classProtectionFilter.Text
		}

		filteredRecords, err := services.SearchSZIRecords(db, uint(user.ID), filters)
		if err != nil {
			// Показать ошибку
			errorDialog := widget.NewModalPopUp(
				widget.NewLabel("Ошибка при поиске записей: "+err.Error()),
				myWindow.Canvas(),
			)
			errorDialog.Show()
			return
		}

		// Обновляем данные в таблице
		updateTableData(filteredRecords)
	}

	// Кнопка сброса фильтров
	resetFilters := func() {
		nameFilter.SetText("")
		typeFilter.SetText("")
		statusFilter.SetText("")
		locationFilter.SetText("")
		manufacturerFilter.SetText("")
		purposeFilter.SetText("")
		deploymentTypeFilter.SetText("")
		classProtectionFilter.SetText("")

		// Загружаем все записи снова
		allRecords, err := services.GetUserRecords(db, uint(user.ID))
		if err != nil {
			// Показать ошибку
			errorDialog := widget.NewModalPopUp(
				widget.NewLabel("Ошибка при загрузке записей: "+err.Error()),
				myWindow.Canvas(),
			)
			errorDialog.Show()
			return
		}

		// Обновляем данные в таблице
		updateTableData(allRecords)
	}

	// Контейнер для фильтров
	filterContainer := container.NewGridWithColumns(9,
		container.NewVBox(widget.NewLabel("Поиск по наименованию:"), nameFilter),
		container.NewVBox(widget.NewLabel("Тип СЗИ:"), typeFilter),
		container.NewVBox(widget.NewLabel("Статус:"), statusFilter),
		container.NewVBox(widget.NewLabel("Место установки:"), locationFilter),
		container.NewVBox(widget.NewLabel("Производитель:"), manufacturerFilter),
		container.NewVBox(widget.NewLabel("Назначение:"), purposeFilter),
		container.NewVBox(widget.NewLabel("Тип развертывания:"), deploymentTypeFilter),
		container.NewVBox(widget.NewLabel("Класс защищенности:"), classProtectionFilter),
		container.NewVBox(
			widget.NewButton("Применить фильтры", applyFilters),
			widget.NewButton("Сбросить фильтры", resetFilters),
		),
	)

	// Кнопки управления
	addBtn := widget.NewButton("Добавить СЗИ", func() {
		ShowAddSziWindow(myApp, uint(user.ID), db)
	})

	editBtn := widget.NewButton("Редактировать", func() {
		dialog := widget.NewModalPopUp(
			widget.NewLabel("Пожалуйста, выберите запись для редактирования, нажав на столбец 'Действия' в нужной строке"),
			myWindow.Canvas(),
		)
		dialog.Show()
	})

	deleteBtn := widget.NewButton("Удалить", func() {
		dialog := widget.NewModalPopUp(
			widget.NewLabel("Пожалуйста, удалите запись, нажав на столбец 'Действия' в нужной строке"),
			myWindow.Canvas(),
		)
		dialog.Show()
	})

	// Кнопки навигации по страницам
	prevPageBtn := widget.NewButton("Предыдущая", func() {
		if currentPage > 0 {
			currentPage--
			updateTableData(allRecords)
		}
	})

	nextPageBtn := widget.NewButton("Следующая", func() {
		totalRecords := len(allRecords)
		if totalRecords == 0 {
			// Если нет записей, не позволяем перейти на следующую страницу
			return
		}
		maxPages := (totalRecords + pageSize - 1) / pageSize
		if currentPage < maxPages-1 {
			currentPage++
			updateTableData(allRecords)
		}
	})

	updatePageInfo()

	// Контейнер для навигации по страницам
	paginationContainer := container.NewHBox(
		prevPageBtn,
		pageInfoLabel,
		nextPageBtn,
	)


	importBtn := widget.NewButton("Импорт из CSV", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err != nil {
				ui.ShowErrorDialog("Ошибка при открытии файла: "+err.Error(), myWindow.Canvas())
				return
			}
			if uri == nil {
				return
			}

			// Получаем путь к файлу
			filePath := uri.URI().String()[7:]

			// Вызываем сервис импорта
			importedRecords, err := csvimport.ImportFromCSV(filePath, uint(user.ID))
			if err != nil {
				errorDialog := widget.NewModalPopUp(
					widget.NewLabel("Ошибка при импорте из CSV: "+err.Error()),
					myWindow.Canvas(),
				)
				errorDialog.Show()
				return
			}

			// Сохраняем импортированные записи в базу данных
			var savedCount int
			for _, record := range importedRecords {
				var existingRecord models.SZIRecord
				result := db.Where("name = ? AND cert_number = ? AND user_id = ?", record.Name, record.CertNumber, record.UserID).First(&existingRecord)

				if result.Error != nil && result.Error.Error() == "record not found" {
					// Записи с такими параметрами не существует, создаём новую (без указания ID, чтобы база данных сгенерировала новый)
					newRecord := models.SZIRecord{
						Name:              record.Name,
						Type:              record.Type,
						CertNumber:        record.CertNumber,
						CertIssueDate:     record.CertIssueDate,
						CertExpiryDate:    record.CertExpiryDate,
						Location:          record.Location,
						Status:            record.Status,
						UserID:            record.UserID,
						Manufacturer:      record.Manufacturer,
						SoftwareVersion:   record.SoftwareVersion,
						Purpose:           record.Purpose,
						DeploymentType:    record.DeploymentType,
						ClassProtection:   record.ClassProtection,
						CreatedAt:         record.CreatedAt,
						UpdatedAt:         record.UpdatedAt,
					}

					err := db.Create(&newRecord).Error
					if err != nil {
						// Показываем ошибку
						errorDialog := widget.NewModalPopUp(
							widget.NewLabel("Ошибка при сохранении записи: "+err.Error()),
							myWindow.Canvas(),
						)
						errorDialog.Show()
						return
					}
					savedCount++
				} else if result.Error == nil {
					existingRecord.Name = record.Name
					existingRecord.Type = record.Type
					existingRecord.CertNumber = record.CertNumber
					existingRecord.CertIssueDate = record.CertIssueDate
					existingRecord.CertExpiryDate = record.CertExpiryDate
					existingRecord.Location = record.Location
					existingRecord.Status = record.Status
					existingRecord.Manufacturer = record.Manufacturer
					existingRecord.SoftwareVersion = record.SoftwareVersion
					existingRecord.Purpose = record.Purpose
					existingRecord.DeploymentType = record.DeploymentType
					existingRecord.ClassProtection = record.ClassProtection
					existingRecord.CreatedAt = record.CreatedAt
					existingRecord.UpdatedAt = record.UpdatedAt

					err := db.Save(&existingRecord).Error
					if err != nil {
						// Показываем ошибку
						errorDialog := widget.NewModalPopUp(
							widget.NewLabel("Ошибка при обновлении записи: "+err.Error()),
							myWindow.Canvas(),
						)
						errorDialog.Show()
						return
					}
					savedCount++
				} else {
					// Другая ошибка при поиске
					errorDialog := widget.NewModalPopUp(
						widget.NewLabel("Ошибка при проверке существующей записи: "+result.Error.Error()),
						myWindow.Canvas(),
					)
					errorDialog.Show()
					return
				}
			}

			// Показываем сообщение об успешном импорте
			var successDialog *widget.PopUp

			successDialog = widget.NewModalPopUp(
				container.NewVBox(
					widget.NewLabel(fmt.Sprintf("Успешно обработано %d записей из CSV", savedCount)),
					widget.NewButton("OK", func() {
						successDialog.Hide()
						// Обновляем окно
						ShowSZIRegistryWindow(myApp, username, db)
						myWindow.Close()
					}),
				),
				myWindow.Canvas(),
			)
			successDialog.Show()
		}, myWindow)
	})

	// Кнопка экспорта в CSV
	exportBtn := widget.NewButton("Экспорт в CSV", func() {
		// Получаем все записи пользователя для экспорта
		records, err := services.GetUserRecords(db, uint(user.ID))
		if err != nil {
			// Показываем ошибку
			errorDialog := widget.NewModalPopUp(
				widget.NewLabel("Ошибка при загрузке данных для экспорта: "+err.Error()),
				myWindow.Canvas(),
			)
			errorDialog.Show()
			return
		}

		// Генерируем имя файла с временной меткой
		filename := fmt.Sprintf("szi_export_%s.csv", time.Now().Format("20060102_150405"))

		// Экспортируем данные
		err = export.ExportToCSV(records, filename)
		if err != nil {
			// Показываем ошибку
			errorDialog := widget.NewModalPopUp(
				widget.NewLabel("Ошибка при экспорте в CSV: "+err.Error()),
				myWindow.Canvas(),
			)
			errorDialog.Show()
			return
		}

		// Показываем сообщение об успешном экспорте
		var exportSuccessDialog *widget.PopUp

		exportSuccessDialog = widget.NewModalPopUp(
			container.NewVBox(
				widget.NewLabel("Данные успешно экспортированы в файл: "+filename),
				widget.NewButton("OK", func() {
					exportSuccessDialog.Hide()
				}),
			),
			myWindow.Canvas(),
		)
		exportSuccessDialog.Show()
	})

	// Контейнер для кнопок управления
	controlButtonsContainer := container.NewHBox(addBtn, editBtn, deleteBtn, importBtn, exportBtn)

	// Обернем таблицу в прокручиваемый контейнер для отображения большего количества строк
	scrollContainer := container.NewVScroll(table)

	// Создаем контейнер с фиксированной высотой для области с таблицей
	tableArea := container.NewVBox(scrollContainer)
	tableArea.Objects[0].(*container.Scroll).SetMinSize(fyne.NewSize(1400, 400))

	// Создаем контейнер с заголовками и таблицей, чтобы они прокручивались вместе
	tableWithHeaders := container.NewVBox(
		headerTable,  // Используем таблицу-заголовок вместо headerContainer
		tableArea,
	)

	// Увеличиваем размер окна для отображения большего количества строк
	myWindow.Resize(fyne.NewSize(1400, 800))

	// Создаем контейнер с правильным расположением элементов
	content := container.NewBorder(
		filterContainer, // верхняя часть - фильтры
		container.NewBorder(nil, paginationContainer, nil, nil, controlButtonsContainer), // нижняя часть - кнопки управления и навигация
		nil, // левая часть - нет
		nil, // правая часть - нет
		tableWithHeaders, // центральная часть - заголовки и таблица
	)

	// Добавляем панель инструментов
	toolbar := container.NewHBox(
		widget.NewButton("Назад", func() {
			myWindow.Close()
			ShowDashboardWindow(myApp, username, db)
		}),
		widget.NewButton("Выход", func() {
			myWindow.Close()
			ShowLoginWindow(myApp, db)
		}),
	)

	// Объединяем основной контент с панелью инструментов
	finalContent := container.NewBorder(toolbar, nil, nil, nil, content)

	myWindow.SetContent(finalContent)

	// Инициализируем данные
	updateTableData(records)

	myWindow.Show()
}