package ui

import (
	"fmt"
	"strings"
	"szi-registry/export"
	"szi-registry/models"
	"szi-registry/services"
	"szi-registry/utils/ui"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

// ShowMainWindow отображает главное окно приложения
func ShowMainWindow(myApp fyne.App, username string, db *gorm.DB) {
	myWindow := myApp.NewWindow("Реестр СЗИ от НСД — " + username)
	myWindow.Resize(fyne.NewSize(800, 500))

	// Получаем ID пользователя по имени
	user, err := services.GetUserByUsername(db, username)
	if err != nil || user == nil {
		// Если не удалось получить пользователя, используем 0
		user = &models.User{ID: 0}
	}

	// Создаем срезы для хранения данных и действий
	var tableData [][]string
	var buttonActions [][2]func()
	var table *widget.Table

	// Загрузка данных из базы данных для конкретного пользователя
	records, err := services.GetUserRecords(db, uint(user.ID))
	if err != nil {
		// В реальном приложении нужно обработать ошибку
		records = []models.SZIRecord{} // Используем пустой массив
	}

	// Функция для обновления данных в таблице
	updateTableData := func(data []models.SZIRecord) {
		tableData = make([][]string, 0)
		buttonActions = make([][2]func(), 0)

		for _, record := range data {
			// Пропускаем записи с пустым названием
			if record.Name == "" {
				continue
			}

			// Определяем статус на основе даты истечения
			status := record.Status

			// Если статус не установлен явно, определяем его по дате истечения
			if status == "" {
				status = ui.GetRecordStatus(record)
			}

			rowData := []string{
				record.Name,
				record.Type,
				record.CertNumber,
				record.CertIssueDate.Format("2006-01-02"),
				record.CertExpiryDate.Format("2006-01-02"),
				status,
				record.Manufacturer,
				record.SoftwareVersion,
				record.Purpose,
				record.DeploymentType,
				record.ClassProtection,
			}

			tableData = append(tableData, rowData)

			// Сохраняем действия для каждой строки
			currentRecord := record // захватываем значение для замыкания
			var confirmDialog *widget.PopUp
			actions := [2]func(){
				func() { // Редактировать
					ShowEditSziWindow(myApp, currentRecord, db)
				},
				func() { // Удалить
					// Создаем диалог подтверждения вручную
					confirmDialog = widget.NewModalPopUp(
						container.NewVBox(
							widget.NewLabel("Вы уверены, что хотите удалить запись '"+currentRecord.Name+"'?"),
							container.NewHBox(
								widget.NewButton("Да", func() {
									err := services.DeleteSziRecord(db, uint(user.ID), currentRecord.ID)
									if err != nil {
										// Показать ошибку
										ui.ShowErrorDialog("Ошибка при удалении записи: "+err.Error(), myWindow.Canvas())
									} else {
										// Обновить окно
										myWindow.Close()
										ShowMainWindow(myApp, username, db)
									}
									confirmDialog.Hide()
								}),
								widget.NewButton("Нет", func() {
									confirmDialog.Hide()
								}),
							),
						),
						myWindow.Canvas(),
					)
					confirmDialog.Show()
				},
			}
			buttonActions = append(buttonActions, actions)
		}
	}

	// Инициализируем данные
	updateTableData(records)

	// Определяем начальные ширины столбцов
	initialColumnWidths := []float32{120, 70, 80, 70, 70, 60, 100, 70, 70, 70, 70, 60}

	// Создаем таблицу с 12 столбцами (11 данных + 1 для действий)
	table = widget.NewTable(
		func() (int, int) {
			return len(tableData), 12 // строки: записи из БД, 12 столбцов
		},
		func() fyne.CanvasObject {
			// Возвращаем метку для обычных ячеек
			label := widget.NewLabel("")
			return label
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			if id.Row < len(tableData) {
				if id.Col < 11 {
					// Обычные данные
					if id.Col < len(tableData[id.Row]) {
						text := tableData[id.Row][id.Col]
						label.SetText(text)
					} else {
						label.SetText("")
					}
				} else if id.Col == 11 {
					// Столбец действий
					label.SetText("...")
					label.Alignment = fyne.TextAlignCenter
				}
			} else {
				label.SetText("")
			}
		})

	// Устанавливаем начальные ширины столбцов
	for i, width := range initialColumnWidths {
		if i < len(initialColumnWidths) {
			table.SetColumnWidth(i, width)
		}
	}

	// Обработка нажатий на таблицу
	table.OnSelected = func(id widget.TableCellID) {
		if id.Col == 11 && id.Row < len(buttonActions) {
			// Показываем диалог выбора действия
			actionMenu := widget.NewPopUpMenu(
				&fyne.Menu{
					Label: "Действия",
					Items: []*fyne.MenuItem{
						{Label: "Редактировать", Action: buttonActions[id.Row][0]},
						{Label: "Удалить", Action: buttonActions[id.Row][1]},
					},
				}, myWindow.Canvas())

			// Показываем меню в позиции курсора
			actionMenu.Show()
		}
	}

	// Заголовки столбцов
	headers := []string{"Наименование СЗИ", "Тип СЗИ", "№ сертификата", "Дата выдачи", "Срок действия", "Статус", "Производитель", "Версия ПО", "Назначение", "Тип развертывания", "Класс защищенности", "Действия"}
	headerContainer := container.NewHBox()
	for _, header := range headers {
		label := widget.NewLabelWithStyle(header, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		label.Alignment = fyne.TextAlignCenter
		box := container.NewVBox(label, widget.NewSeparator())
		headerContainer.Add(box)
	}

	// Создаем поля для фильтрации
	nameFilter := widget.NewEntry()
	nameFilter.SetPlaceHolder("Поиск по наименованию...")
	typeFilter := widget.NewSelectEntry([]string{"СКЗИ", "ОС", "СЗИ КС", "СЗИ СКЗИ", "Другое"})
	typeFilter.SetPlaceHolder("Тип СЗИ...")
	statusFilter := widget.NewSelectEntry([]string{"Актуально", "Требует внимания", "Просрочено"})
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
			ui.ShowErrorDialog("Ошибка при поиске записей: "+err.Error(), myWindow.Canvas())
			return
		}

		// Обновляем данные в таблице
		updateTableData(filteredRecords)
		table.Refresh()
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
			ui.ShowErrorDialog("Ошибка при загрузке записей: "+err.Error(), myWindow.Canvas())
			return
		}

		// Обновляем данные в таблице
		updateTableData(allRecords)
		table.Refresh()
	}

	// Добавляем кнопки для межпользовательского взаимодействия
	shareBtn := widget.NewButton("Поделиться записью", func() {
		// Открываем диалог для выбора записи и пользователя для предоставления доступа
		dialog := createShareDialog(myWindow, uint(user.ID), db)
		dialog.Show()
	})

	// Кнопка импорта из CSV
	importBtn := widget.NewButton("Импорт из CSV", func() {
		// Диалог для выбора CSV файла
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog := widget.NewModalPopUp(
					widget.NewLabel("Ошибка при открытии файла: "+err.Error()),
					myWindow.Canvas(),
				)
				dialog.Show()
				return
			}

			if reader == nil {
				return // Пользователь отменил выбор файла
			}

			// Получаем путь к файлу
			filePath := reader.URI().String()

			// Вызываем сервис импорта
			importedCount, err := services.ImportSziRecordsFromCSV(db, uint(user.ID), filePath)
			if err != nil {
				dialog := widget.NewModalPopUp(
					widget.NewLabel("Ошибка при импорте из CSV: "+err.Error()),
					myWindow.Canvas(),
				)
				dialog.Show()
			} else {
				// Обновляем данные в таблице
				records, err := services.GetUserRecords(db, uint(user.ID))
				if err != nil {
					// В реальном приложении нужно обработать ошибку
					records = []models.SZIRecord{} // Используем пустой массив
				}

				// Обновляем данные таблицы
				updateTableData(records)
				table.Refresh()

				// Показываем сообщение об успешном импорте
				dialog := widget.NewModalPopUp(
					widget.NewLabel(fmt.Sprintf("Успешно импортировано %d записей из CSV", importedCount)),
					myWindow.Canvas(),
				)
				dialog.Show()
			}
		}, myWindow)

		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
		fileDialog.Show()
	})

	// Кнопка экспорта в CSV
	exportBtn := widget.NewButton("Экспорт в CSV", func() {
		// Получаем все записи пользователя для экспорта
		records, err := services.GetUserRecords(db, uint(user.ID))
		if err != nil {
			dialog := widget.NewModalPopUp(
				widget.NewLabel("Ошибка при загрузке данных для экспорта: "+err.Error()),
				myWindow.Canvas(),
			)
			dialog.Show()
			return
		}

		// Генерируем имя файла с временной меткой
		filename := fmt.Sprintf("szi_export_%s.csv", time.Now().Format("20060102_150405"))

		// Экспортируем данные
		err = export.ExportToCSV(records, filename)
		if err != nil {
			dialog := widget.NewModalPopUp(
				widget.NewLabel("Ошибка при экспорте в CSV: "+err.Error()),
				myWindow.Canvas(),
			)
			dialog.Show()
		} else {
			dialog := widget.NewModalPopUp(
				widget.NewLabel("Данные успешно экспортированы в файл: "+filename),
				myWindow.Canvas(),
			)
			dialog.Show()
		}
	})

	addBtn := widget.NewButton("Добавить СЗИ", func() {
		user, err := services.GetUserByUsername(db, username)
		if err != nil || user == nil {
			// Если не удалось получить пользователя, показываем ошибку
			dialog := widget.NewModalPopUp(
				widget.NewLabel("Не удалось получить информацию о пользователе"),
				myWindow.Canvas(),
			)
			dialog.Show()
			return
		}
		ShowAddSziWindow(myApp, uint(user.ID), db)
	})

	logoutBtn := widget.NewButton("Выход", func() {
		myWindow.Close()
		ShowLoginWindow(myApp, db)
	})

	// Кнопки фильтров
	applyFiltersBtn := widget.NewButton("Применить фильтры", applyFilters)
	resetFiltersBtn := widget.NewButton("Сбросить фильтры", resetFilters)

	// Контейнер для фильтров
	filterContainer := container.NewGridWithColumns(6,
		container.NewVBox(widget.NewLabel("Поиск по наименованию:"), nameFilter),
		container.NewVBox(widget.NewLabel("Тип СЗИ:"), typeFilter),
		container.NewVBox(widget.NewLabel("Статус:"), statusFilter),
		container.NewVBox(widget.NewLabel("Место установки:"), locationFilter),
		container.NewVBox(widget.NewLabel("Производитель:"), manufacturerFilter),
		container.NewVBox(widget.NewLabel("Назначение:"), purposeFilter),
	)

	// Создаем контейнер для верхней панели с основными кнопками и фильтрами
	topToolbar := container.NewHBox(
		logoutBtn,
		addBtn,
		applyFiltersBtn,
		resetFiltersBtn,
	)

	// Создаем контейнер с заголовками и таблицей, чтобы они прокручивались вместе
	tableWithHeaders := container.NewVBox(
		headerContainer,
		table,
	)

	// Добавляем нижнюю панель с дополнительными кнопками
	bottomPanel := container.NewHBox(
		shareBtn,
		importBtn,
		exportBtn,
		widget.NewButton("Статистика", func() {
			showStatisticsWindow(myApp, uint(user.ID), db)
		}),
	)

	content := container.NewBorder(
		filterContainer,  // верхняя часть - фильтры
		bottomPanel,      // нижняя часть - дополнительные кнопки
		nil,              // левая часть - нет
		nil,              // правая часть - нет
		tableWithHeaders, // центральная часть - заголовки и таблица
	)
	// Объединяем основной контент с верхней панелью инструментов
	finalContent := container.NewBorder(topToolbar, nil, nil, nil, content)

	// Обернем в прокрутку для больших таблиц
	scrollableContent := container.NewScroll(finalContent)
	scrollableContent.SetMinSize(fyne.NewSize(700, 400))

	myWindow.SetContent(scrollableContent)
	myWindow.Show()
}

// createShareDialog создает диалог для предоставления доступа к записи другому пользователю
func createShareDialog(parentWindow fyne.Window, currentUserID uint, db *gorm.DB) *widget.PopUp {
	// Создаем элементы интерфейса для выбора записи и пользователя
	recordSelector := widget.NewEntry()
	recordSelector.SetPlaceHolder("Введите название СЗИ")
	userSelector := widget.NewSelectEntry([]string{})
	permissionSelector := widget.NewSelect([]string{"r", "rw"}, func(string) {})

	// Создаем контейнер для предложений записей
	suggestionsContainer := container.NewVBox()

	// Создаем диалог
	var dialog *widget.PopUp

	// Создаем списки записей и пользователей
	var allRecords []models.SZIRecord
	var filteredRecords []models.SZIRecord

	// Загружаем доступные записи и пользователей
	loadData := func() {
		// Загружаем записи пользователя
		records, err := services.GetUserRecords(db, currentUserID)
		if err != nil {
			// Обработка ошибки
			fmt.Printf("Ошибка при загрузке записей: %v\n", err)
			return
		}

		allRecords = records
		filteredRecords = records

		// Загружаем всех пользователей, кроме текущего
		allUsers, err := services.GetAllUsers(db)
		if err != nil {
			// Обработка ошибки
			fmt.Printf("Ошибка при загрузке пользователей: %v\n", err)
			return
		}

		var userOptions []string
		for _, user := range allUsers {
			if user.ID != currentUserID { // Исключаем текущего пользователя
				userOptions = append(userOptions, user.Username)
			}
		}
		userSelector.SetOptions(userOptions)
	}

	// Функция для обновления списка записей на основе введенного текста
	updateSuggestions := func() {
		// Очищаем предыдущие предложения
		suggestionsContainer.Objects = nil

		// Фильтруем записи по введенному тексту
		inputText := recordSelector.Text
		filteredRecords = []models.SZIRecord{}

		for _, record := range allRecords {
			if strings.Contains(strings.ToLower(record.Name), strings.ToLower(inputText)) {
				filteredRecords = append(filteredRecords, record)
			}
		}

		// Создаем кнопки для каждой подходящей записи
		for _, record := range filteredRecords {
			btn := widget.NewButton(record.Name, func() {
				// При выборе записи, устанавливаем её название в поле ввода
				recordSelector.SetText(record.Name)
				// Скрываем список предложений
				suggestionsContainer.Hide()
			})
			suggestionsContainer.Add(btn)
		}

		// Обновляем интерфейс
		if len(filteredRecords) > 0 {
			suggestionsContainer.Show()
		} else {
			suggestionsContainer.Hide()
		}
		suggestionsContainer.Refresh()
	}

	loadData()

	// Обработчик изменения текста в поле ввода
	recordSelector.OnChanged = func(text string) {
		updateSuggestions()
	}

	// Кнопка подтверждения
	confirmBtn := widget.NewButton("Поделиться", func() {
		// Функция предоставления доступа больше не поддерживается
		ui.ShowErrorDialog("Функция предоставления доступа больше не поддерживается", parentWindow.Canvas())
	})

	// Кнопка отмены
	cancelBtn := widget.NewButton("Отмена", func() {
		dialog.Hide()
	})

	dialogContent := container.NewVBox(
		widget.NewLabel("Название СЗИ:"),
		recordSelector,
		suggestionsContainer,
		widget.NewLabel("Выберите пользователя:"),
		userSelector,
		widget.NewLabel("Выберите права доступа:"),
		permissionSelector,
		container.NewHBox(confirmBtn, cancelBtn),
	)

	dialog = widget.NewModalPopUp(dialogContent, parentWindow.Canvas())
	return dialog
}
