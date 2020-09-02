heat dir ..\web\build -nologo -cg AppFiles0 -gg -g1 -srd -sfrag -template fragment -dr APPDIR0 -var var.SourceDir0 -out AppFiles0.wxs
candle -dSourceDir0=..\web\build AppFiles0.wxs LicenseAgreementDlg_HK.wxs WixUI_HK.wxs product.wxs
light -ext WixUIExtension -ext WixUtilExtension -sacl -spdb  -out ..\wix\prom2lyrid.msi AppFiles0.wixobj LicenseAgreementDlg_HK.wixobj WixUI_HK.wixobj product.wixobj
