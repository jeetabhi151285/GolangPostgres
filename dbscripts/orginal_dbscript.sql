USE [Maersk]
GO
/****** Object:  Table [dbo].[xPreBaplieDetails]    Script Date: 11/01/2020 23:49:53 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
SET ANSI_PADDING ON
GO
CREATE TABLE [dbo].[xPreBaplieDetails](
	[vessel] [varchar] (50) NULL,
	[callSign] [varchar](50) NULL,
	[voyage] [varchar](200) NULL,
	[carrier] [varchar](200) NULL,
	[operation] [varchar](100) NULL,
	[creationDate] [varchar](100) NULL,
	[containerDetails] json NOT NULL
) ON [PRIMARY]
GO
SET ANSI_PADDING OFF
GO
/****** Object:  Table [dbo].[xSampleTestTable]    Script Date: 13/01/2020 11:06:53 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
SET ANSI_PADDING ON
GO
CREATE TABLE xSampleTestTable(
	vessel varchar  NULL,
	callSign varchar NULL,
	voyage varchar NULL,
	carrier varchar NULL,
	operation varchar NULL,
	creationDate varchar NULL,
	containerDetails json NOT NULL
)
GO
SET ANSI_PADDING OFF
GO