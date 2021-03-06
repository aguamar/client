// @flow
import React, {Component} from 'react'
import {Box, Box2, Icon, InfoNote, Text, Button, Input, PlatformIcon} from '../../common-adapters'
import {globalStyles, globalColors, globalMargins, desktopStyles, collapseStyles} from '../../styles'
import {platformText} from './shared'
import type {PlatformsExpandedType} from '../../constants/types/more'
import type {Props} from '.'
import openURL from '../../util/open-url'

function UsernameTips({platform}: {platform: PlatformsExpandedType}) {
  if (platform === 'hackernews') {
    return (
      <Box style={styleYellowBanner}>
        <Text backgroundMode="Information" type="BodySmallSemibold">
          &bull; You must have karma &ge; 2<br />
          &bull; You must enter your uSeRName with exact case
        </Text>
      </Box>
    )
  }

  if (platform === 'facebook') {
    return (
      <InfoNote containerStyle={styleInfoNote}>
        <Box2 direction="vertical" style={{textAlign: 'center'}}>
          <Text type="BodySmall" style={{textAlign: 'center'}}>
            You can find your Facebook username at
          </Text>
          <Box2 direction="horizontal">
            <Text type="BodySmallSecondaryLink" onClick={() => openURL('http://www.facebook.com/settings')}>
              http://www.facebook.com/settings
            </Text>
            <Text type="BodySmall">.</Text>
          </Box2>
        </Box2>
      </InfoNote>
    )
  }

  return null
}

type State = {
  username: string,
}

function customError(error: string, code: ?number) {
  return (
    <Text style={styleErrorBannerText} type="BodySemibold">
      {error}
    </Text>
  )
}

class PrivateEnterUsernameRender extends Component<Props, State> {
  state: State

  constructor(props: Props) {
    super(props)
    this.state = {
      username: '',
    }
  }

  handleUsernameChange(username: string) {
    if (this.props.onUsernameChange) {
      this.props.onUsernameChange(username)
    }
    this.setState({username})
  }

  handleContinue() {
    this.props.onContinue(this.state.username)
  }

  render() {
    const {headerText, floatingLabelText, hintText} = platformText[this.props.platform]

    return (
      <Box style={styleContainer}>
        <Icon style={styleClose} type="iconfont-close" onClick={this.props.onCancel} />
        {this.props.errorText && (
          <Box style={styleErrorBanner}>{customError(this.props.errorText, this.props.errorCode)}</Box>
        )}
        <Text type="Header" style={{marginBottom: globalMargins.medium}}>
          {headerText}
        </Text>
        <PlatformIcon
          platform={this.props.platform}
          overlay={'icon-proof-unfinished'}
          overlayColor={globalColors.grey}
        />
        <Input
          autoFocus={true}
          style={styleInput}
          floatingHintTextOverride={floatingLabelText}
          hintText={hintText}
          value={this.state.username}
          onChangeText={username => this.handleUsernameChange(username)}
          onEnterKeyDown={() => this.handleContinue()}
        />
        <UsernameTips platform={this.props.platform} />
        <Box style={{...globalStyles.flexBoxRow, marginTop: 32}}>
          <Button
            type="Secondary"
            onClick={this.props.onCancel}
            label="Cancel"
            style={{marginRight: globalMargins.tiny}}
          />
          <Button
            type="Primary"
            disabled={!this.props.canContinue}
            onClick={() => this.handleContinue()}
            label="Continue"
          />
        </Box>
      </Box>
    )
  }
}

// FIXME: this is the old way (#styles)

const styleErrorBanner = {
  ...globalStyles.flexBoxColumn,
  justifyContent: 'center',
  position: 'absolute',
  alignItems: 'center',
  top: 0,
  left: 0,
  right: 0,
  zIndex: 1,
  minHeight: globalMargins.large,
  backgroundColor: globalColors.red,
}

const styleErrorBannerText = {
  color: globalColors.white,
}

const styleContainer = {
  ...globalStyles.flexBoxColumn,
  position: 'relative',
  flex: 1,
  alignItems: 'center',
  justifyContent: 'center',
}

const styleClose = collapseStyles([
  desktopStyles.clickable,
  {
    position: 'absolute',
    right: 16,
    top: 16,
  },
])

const styleInput = {
  alignSelf: 'center',
  marginTop: globalMargins.small,
  marginBottom: 0,
  width: 460,
}

const styleYellowBanner = {
  ...globalStyles.flexBoxColumn,
  alignItems: 'center',
  backgroundColor: globalColors.yellow,
  borderRadius: 3,
  marginTop: globalMargins.small,
  marginBottom: -globalMargins.tiny,
  minWidth: 460,
  paddingTop: globalMargins.xsmall,
  paddingBottom: globalMargins.xsmall,
  paddingLeft: globalMargins.small,
  paddingRight: globalMargins.small,
}

const styleInfoNote = {
  marginBottom: globalMargins.tiny,
  marginTop: globalMargins.medium,
}

export default PrivateEnterUsernameRender
