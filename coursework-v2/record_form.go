package main

import (
  "fmt"

  "szi-registry/coursework-v2/internal/application"
  "szi-registry/coursework-v2/internal/domain"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/widget"
)

func showRecordForm(myApp fyne.App, app application.App, userID uint, record *domain.SZIRecord, onClosed func()) {
  editing := record != nil
  title := "Добавить средство защиты"
  if editing {
    title = "Редактировать средство защиты"
  }
  window := myApp.NewWindow(title)
  window.Resize(fyne.NewSize(720, 700))
  name := widget.NewEntry()
  typeSelect := widget.NewSelect(domain.SZITypeValues, nil)
  protection := widget.NewSelect(domain.ProtectionClassValues, nil)
  trust := widget.NewSelect(domain.TrustLevelValues, nil)
  ac := widget.NewSelect(domain.ACClassValues, nil)
  vendor := widget.NewEntry()
  version := widget.NewEntry()
  certificate := widget.NewEntry()
  scheme := widget.NewSelect(domain.CertificationSchemeValues, nil)
  issue := widget.NewEntry()
  expiry := widget.NewEntry()
  location := widget.NewEntry()
  responsible := widget.NewEntry()
  notes := widget.NewMultiLineEntry()
  issue.SetPlaceHolder("ГГГГ-ММ-ДД")
  expiry.SetPlaceHolder("ГГГГ-ММ-ДД")
  if editing {
    name.SetText(record.Name)
    typeSelect.SetSelected(record.SZIType)
    protection.SetSelected(record.ProtectionClass)
    trust.SetSelected(record.TrustLevel)
    ac.SetSelected(record.ACClass)
    vendor.SetText(record.Vendor)
    version.SetText(record.Version)
    certificate.SetText(record.CertificateNumber)
    scheme.SetSelected(record.CertificationScheme)
    issue.SetText(record.IssueDate.Format("2006-01-02"))
    expiry.SetText(record.ExpiryDate.Format("2006-01-02"))
    location.SetText(record.InstallLocation)
    responsible.SetText(record.ResponsiblePerson)
    notes.SetText(record.Notes)
  }
  form := widget.NewForm(widget.NewFormItem("Наименование СЗИ", name), widget.NewFormItem("Тип СЗИ", typeSelect), widget.NewFormItem("Класс защищённости СВТ", protection), widget.NewFormItem("Уровень доверия", trust), widget.NewFormItem("Класс АС", ac), widget.NewFormItem("Заявитель / разработчик", vendor), widget.NewFormItem("Версия ПО / исполнение", version), widget.NewFormItem("Номер сертификата", certificate), widget.NewFormItem("Схема сертификации", scheme), widget.NewFormItem("Дата выдачи", issue), widget.NewFormItem("Срок действия", expiry), widget.NewFormItem("Место установки", location), widget.NewFormItem("Ответственный сотрудник", responsible), widget.NewFormItem("Примечания", notes))
  form.OnSubmit = func() {
    if name.Text == ""  typeSelect.Selected == ""  certificate.Text == "" {
      dialog.ShowError(fmt.Errorf("заполните наименование, тип СЗИ и номер сертификата"), window)
      return
    }
    issueDate, err := domain.ParseDate(issue.Text)
    if err != nil {
      dialog.ShowError(fmt.Errorf("дата выдачи: %v", err), window)
      return
    }
    expiryDate, err := domain.ParseDate(expiry.Text)
    if err != nil {
      dialog.ShowError(fmt.Errorf("срок действия: %v", err), window)
      return
    }
    if record == nil {
      record = &domain.SZIRecord{OwnerID: userID}
    }
    record.Name = name.Text
    record.SZIType = typeSelect.Selected
    record.ProtectionClass = protection.Selected
    record.TrustLevel = trust.Selected
    record.ACClass = ac.Selected
    record.Vendor = vendor.Text
    record.Version = version.Text
    record.CertificateNumber = certificate.Text
    record.CertificationScheme = scheme.Selected
    record.IssueDate = issueDate
    record.ExpiryDate = expiryDate
    record.InstallLocation = location.Text
    record.ResponsiblePerson = responsible.Text
    record.Notes = notes.Text
    if editing {
      err = app.Update(*record)
    } else {
      err = app.Create(*record)
    }
    if err != nil {
      dialog.ShowError(fmt.Errorf("не удалось сохранить запись: %v", err), window)
      return
    }
    window.Close()
  }
  closeButton := widget.NewButton("Закрыть", window.Close)
  closeButton.Importance = widget.HighImportance
  saveButton := widget.NewButton("Сохранить", func() { form.OnSubmit() })
  saveButton.Importance = widget.HighImportance
  window.SetContent(container.NewBorder(nil, container.NewHBox(closeButton, saveButton), nil, nil, container.NewVScroll(form)))
  window.SetOnClosed(onClosed)
  window.Show()
}