.\WiX\heat.exe dir ..\web\build -nologo -cg AppFiles0 -gg -g1 -srd -sfrag -template fragment -dr APPDIR0 -var var.SourceDir0 -out AppFiles0.wxs
.\WiX\candle.exe -ext WixUIExtension -ext WixUtilExtension -dSourceDir0=..\web\build -dVersion=1.0.0 AppFiles0.wxs LicenseAgreementDlg_HK.wxs WixUI_HK.wxs product.wxs
.\WiX\light.exe -ext WixUIExtension -ext WixUtilExtension -sacl -spdb  -out ..\wix\prom2lyrid.msi AppFiles0.wixobj LicenseAgreementDlg_HK.wixobj WixUI_HK.wixobj product.wixobj
